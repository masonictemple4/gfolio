package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/masonictemple4/masonictempl/db"
	"github.com/masonictemple4/masonictempl/db/models"
	"github.com/masonictemple4/masonictempl/internal/dtos"
	"github.com/masonictemple4/masonictempl/internal/filestore"
	"github.com/masonictemple4/masonictempl/internal/parser"
	"github.com/masonictemple4/masonictempl/services"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var blogCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new blog.",
	Long: `Create a new blog. For example:
masonictempl blog create <file path>`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		blogsCmd.PersistentPreRun(blogsCmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		pubDir, err := cmd.PersistentFlags().GetString("pub")
		if err != nil {
			log.Fatal(err)
		}

		if err := createBlog(cmd.Context(), args[0], pubDir, cmd.Flags()); err != nil {
			log.Fatal(err)
		}

	},
}

func createBlog(ctx context.Context, path, pubRoot string, flags *pflag.FlagSet) error {
	// What if instead of passing a path to parser here and eventually would have to be
	// the filestore too if i read file here and pass bytes to the parsr and writer.
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	blogDb := db.NewPostgresGCPProxy(db.WithDSN())
	localService := services.NewBlogService(blogDb)

	var result dtos.BlogInput
	err = parser.ParseFile(path, &result)
	if err != nil {
		return err
	}

	fmt.Printf("\nThe parsed result is:  %+v\n", result)

	post := &models.Blog{
		State: models.BlogStateDraft,
	}

	err = post.FromBlogInput(blogDb, &result)
	if err != nil {
		return err
	}

	// generate a slug because now we have a title.
	post.Slug = post.GenerateSlug("")

	fileHandler, err := filestore.NewInternalStore(pubRoot)
	if err != nil {
		return err
	}

	post.Docpath, err = post.GenerateDocPath(filestore.GetRootPath(fileHandler))
	if err != nil {
		return err
	}

	fmt.Println("[After parsed] The post ID is: ", post.ID)
	fmt.Println("[After parsed] The post docpath is: ", post.Docpath)

	written, err := fileHandler.Write(ctx, post.Docpath, data)
	if err != nil || len(data) != int(written) {
		return err
	}

	updateBody := map[string]any{"contenturl": post.GenerateContentUrl(), "docpath": post.Docpath, "state": models.BlogStatePublished, "slug": post.Slug}
	err = localService.Store.UpdateFromMap(blogDb, post, updateBody, int(post.ID))
	if err != nil {
		return err
	}

	return nil
}

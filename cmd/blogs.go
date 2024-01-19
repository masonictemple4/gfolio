package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/masonictemple4/masonictempl/db/models"
	"github.com/masonictemple4/masonictempl/internal/dtos"
	"github.com/masonictemple4/masonictempl/internal/filestore"
	"github.com/masonictemple4/masonictempl/internal/parser"
	"github.com/masonictemple4/masonictempl/services"
)

var blogsCmd = &cobra.Command{
	Use:   "blog",
	Short: "Manage blogs within the masonictempl project.",
	Long: `Manage blogs within the masonictempl project. For example:
You canuse this command to list, create, update, and delete blogs.
The default will be to list all blogs.
masonictempl blogs [command]`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var blogCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new blog.",
	Long: `Create a new blog. For example:
masonictempl blog create <file path>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		pubDir, err := cmd.Flags().GetString("pub")
		if err != nil {
			log.Fatal(err)
		}

		if err := createBlog(cmd.Context(), args[0], pubDir, cmd.Flags()); err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(blogsCmd)
	blogsCmd.AddCommand(blogCreateCmd)
	blogsCmd.PersistentFlags().String("pub", os.Getenv("ASSET_DIR"), "name of your public/static file directory.")
}

func createBlog(ctx context.Context, path, pubRoot string, flags *pflag.FlagSet) error {
	// What if instead of passing a path to parser here and eventually would have to be
	// the filestore too if i read file here and pass bytes to the parsr and writer.
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	localService := services.NewBlogService()
	blogDb := localService.Store.DB()

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

func updateBlog(ctx context.Context, bid int, localstore, path string) error {

	localService := services.NewBlogService()
	blogDb := localService.Store.DB()

	var blog models.Blog
	err := localService.Store.FindByID(blogDb, &blog, bid)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var result dtos.BlogInput
	err = parser.ParseFile(path, &result)
	if err != nil {
		return err
	}

	fmt.Printf("\nThe parsed result is:  %+v\n", result)

	err = blog.FromBlogInput(blogDb, &result)
	if err != nil {
		return err
	}

	// generate a slug because now we have a title.
	blog.Slug = blog.GenerateSlug("")

	fileHandler, err := filestore.NewInternalStore(localstore)
	if err != nil {
		return err
	}

	// TODO: Should probably delete here to be safe we're not
	// gathering too unused files.
	blog.Docpath, err = blog.GenerateDocPath(filestore.GetRootPath(fileHandler))
	if err != nil {
		return err
	}

	written, err := fileHandler.Write(ctx, blog.Docpath, data)
	if err != nil || len(data) != int(written) {
		return err
	}

	updateBody := map[string]any{
		"contenturl": blog.GenerateContentUrl(),
		"docpath":    blog.Docpath,
		"state":      models.BlogStatePublished,
		"slug":       blog.Slug,
	}
	err = localService.Store.UpdateFromMap(blogDb, &blog, updateBody, int(blog.ID))
	if err != nil {
		return err
	}

	return nil

}

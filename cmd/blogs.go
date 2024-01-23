package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		rootCmd.PersistentPreRun(rootCmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(blogsCmd)
	blogsCmd.AddCommand(blogCreateCmd)
	blogsCmd.AddCommand(blogsListCmd)
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

package cmd

import (
	"context"
	"encoding/json"

	"github.com/masonictemple4/masonictempl/db"
	"github.com/masonictemple4/masonictempl/services"
	"github.com/spf13/cobra"
)

var blogsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all blogs",
	Long:  `List all blogs`,
	Run: func(cmd *cobra.Command, args []string) {
		listBlogs(cmd.Context())
	},
}

func listBlogs(ctx context.Context) {
	bDb := db.NewPostgresGCPProxy()

	bServ := services.NewBlogService(bDb)
	blogs := bServ.List(ctx)

	if len(blogs) == 0 {
		println("No blogs found")
		return
	}

	println("Blogs: ")
	blogJson, _ := json.MarshalIndent(blogs, "", "  ")
	println(string(blogJson))
}

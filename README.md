# Masonictempl
Slowly porting my old react and next.js portoflio/blog to a Go and Htmx solution. 

#### [CLI Docs](./docs/masonictempl.md)

#### Local storage pathing.
- ENV var `ASSET_DIR` (default: **/project/root/assets**: This is the path to the directory that stores
all of your public static files such as; `css`, `js`, `images & icons`. It can also be used
to serve other files if you don't wish to upload them to a cloud service like S3 or GCP Storage buckets. 
Think of this similar to the `public` folder in a standard js project.

- ENV var `BLOG_ROOT` (default: **blogs**): This is the blog root dir path. This is what get's used 
in either the local or remote storage systems.


- Blog storage pattern is like so: `$BLOG_ROOT/YYYY/MM/DD/blog-slug.md`

import Markdoc from "@markdoc/markdoc";
import { readdirSync } from "node:fs";
import { join } from "node:path";
import yaml from "js-yaml";

interface Post {
  slug: string;
  title: string;
  date: string;
  content: string;
}

// Load all posts
async function loadPosts(): Promise<Post[]> {
  const postsDir = join(import.meta.dir, "posts");
  const files = readdirSync(postsDir).filter((f) => f.endsWith(".md"));

  const posts: Post[] = [];

  for (const file of files) {
    const filePath = join(postsDir, file);
    const source = await Bun.file(filePath).text();

    const ast = Markdoc.parse(source);
    const content = Markdoc.transform(ast);
    const frontmatter = ast.attributes.frontmatter
      ? yaml.load(ast.attributes.frontmatter)
      : {};
    const html = Markdoc.renderers.html(content);

    posts.push({
      slug: file.replace(".md", ""),
      title: frontmatter.title,
      date: frontmatter.date || "Unknown",
      content: html,
    });
  }

  // Sort by date, newest first
  return posts.sort(
    (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime(),
  );
}

const posts = await loadPosts();

const server = Bun.serve({
  async fetch(req) {
    const url = new URL(req.url);

    // Home page - list all posts
    if (url.pathname === "/") {
      const postList = posts
        .map(
          (post) => `
        <article>
          <h2><a href="/posts/${post.slug}">${post.title}</a></h2>
          <time>${post.date}</time>
        </article>
      `,
        )
        .join("");

      return new Response(
        `
        <!DOCTYPE html>
        <html>
          <head>
            <meta charset="utf-8">
            <title>Noah Costello's Blog</title>
            <style>
              body { max-width: 800px; margin: 0 auto; padding: 20px; font-family: system-ui; }
              article { margin: 40px 0; }
              time { color: #666; font-size: 0.9em; }
            </style>
          </head>
          <body>
            <h1>Noah Costello</h1>
            ${postList}
          </body>
        </html>
      `,
        {
          headers: { "Content-Type": "text/html" },
        },
      );
    }

    // Individual post pages
    if (url.pathname.startsWith("/posts/")) {
      const slug = url.pathname.replace("/posts/", "");
      const post = posts.find((p) => p.slug === slug);

      if (!post) {
        return new Response("Post not found", { status: 404 });
      }

      return new Response(
        `
        <!DOCTYPE html>
        <html>
          <head>
            <meta charset="utf-8">
            <title>${post.title} - Noah Costello's Blog</title>
            <style>
              body { max-width: 800px; margin: 0 auto; padding: 20px; font-family: system-ui; line-height: 1.6; }
              img { max-width: 100%; height: auto; }
              time { color: #666; font-size: 0.9em; }
              a { color: #0066cc; }
            </style>
          </head>
          <body>
            <nav><a href="/">← Back to home</a></nav>
            <time>${post.date}</time>
            ${post.content}
          </body>
        </html>
      `,
        {
          headers: { "Content-Type": "text/html" },
        },
      );
    }

    return new Response("Not found", { status: 404 });
  },
});

console.log(`Server running at ${server.url}`);

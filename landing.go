package main

import "html/template"

var PostTmpl = template.Must(template.New("post").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>{{.Title}}</title>
	<style>
		body { font-family: sans-serif; max-width: 800px; margin: 50px auto; padding: 0 20px; }
		a { color: #0066cc; text-decoration: none; }
		a:hover { text-decoration: underline; }
		.back { margin-bottom: 20px; }
	</style>
</head>
<body>
	<div class="back"><a href="/blog">← Back to all posts</a></div>
	{{.Content}}
</body>
</html>
`))

var ListTmpl = template.Must(template.New("list").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Blog</title>
	<style>
		body { font-family: sans-serif; max-width: 800px; margin: 50px auto; padding: 0 20px; }
		h1 { color: #333; }
		ul { list-style: none; padding: 0; }
		li { margin: 10px 0; }
		a { color: #0066cc; text-decoration: none; }
		a:hover { text-decoration: underline; }
	</style>
</head>
<body>
	<p>
	Hi, I'm Noah, I'm a Software Engineer. I created this site so that I could start writing blogs.
	I plan on writing about my hobbies (cooking, bikes, books, videogames) and programming or computer
	science topics. I hope you enjoy and please provide feedback so I can become a better writer.
	<br>
	<br>
	<strong>Future topic ideas:</strong> bike stuff, volunteering, distributed systems course labs
	</p>	

	<h1>Coding Stuff</h1>
	<ul>
		<li><a href="https://gitlab.com/nmcostello/blog">Blog Site Code</a></li>
		<li><a href="https://gitlab.com/nmcostello/dist-system-labs">Distribute Systems Course Lab Code</a></li>
	</ul>
	
	<h1>Blog Posts</h1>
	<ul>
	{{range .}}
		<li><a href="/blog/{{.URL}}">{{.Date}} {{.Title}}</a></li>
	{{end}}
	</ul>
</body>
</html>
`))

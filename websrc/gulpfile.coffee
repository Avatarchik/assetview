gulp       = require 'gulp'
coffee     = require 'gulp-coffee'
concat     = require 'gulp-concat'
plumber    = require 'gulp-plumber'
sass       = require 'gulp-sass'
sourcemaps = require 'gulp-sourcemaps'
notify     = require 'gulp-notify'


files =
	coffee: './coffee/**/*.coffee'
	sass  : './sass/**/*.sass'
	html  : './html/**/*.html'


gulp.task 'js', ->
	gulp.src files.coffee
		.on 'error', ()->{}
		.pipe plumber
      errorHandler: notify.onError("Error: <%= error.message %>")
		.pipe sourcemaps.init
			loadMaps: true
		.pipe coffee
			bare: true
		.pipe concat 'app.js'
		.pipe sourcemaps.write '.',
			addComment: true
			sourceRoot: '/src'
		.pipe gulp.dest '../web/'


gulp.task 'css', ->
	gulp.src files.sass
		.pipe plumber
			errorHandler: notify.onError("Error: <%= error.message %>")
		.pipe sass
			indentedSyntax: true
		.pipe concat 'app.css'
		.pipe gulp.dest '../web'

gulp.task 'html', ->
	gulp.src files.html
		.pipe gulp.dest '../web'

gulp.task 'watch', ['build'], ->
	gulp.watch files.coffee, ['js']
	gulp.watch files.sass, ['css']
	gulp.watch files.html, ['html']

gulp.task 'build', ['js', 'css', 'html']
gulp.task 'default', ['build']

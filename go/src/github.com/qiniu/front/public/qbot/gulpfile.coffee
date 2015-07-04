"use strict"
gulp = require 'gulp'
gutil = require 'gulp-util'
path = require 'path'
plumber = require 'gulp-plumber'

less = require 'gulp-less'
coffee = require 'gulp-coffee'
concat = require 'gulp-concat'
minifier = require 'gulp-minifier'

browserSync = require 'browser-sync'
reload = browserSync.reload
runSequence = require 'run-sequence'
    .use gulp

errHandler =  (err) ->
  gutil.beep()
  gutil.log err
  this.emit 'end'

coffeeOption =
  bare: true

minifierOption =
  minify: true
  collapseWhitespace: true
  conservativeCollapse: true
  minifyJS: true
  minifyCSS: true

IS_PROD = false

lessFiles = [
  'css/main.less'
]

gulp.task 'build-less', ->
  gulp.src lessFiles
    .pipe plumber
      errorHandler: errHandler
    .pipe less()
    .pipe concat 'main.min.css'
    .pipe(if IS_PROD then minifier minifierOption else gutil.noop())
    .pipe gulp.dest ''
    .pipe reload {stream: true}

coffeeFiles = [
  'coffee/*.coffee'
]

gulp.task 'coffee', ->
  gulp.src coffeeFiles
    .pipe plumber()
    .pipe coffee(coffeeOption).on('error', gutil.log)
    .pipe gulp.dest 'js/'

jsFiles = [
  'js/*.js'
]

gulp.task 'build-js', ->
  gulp.src jsFiles
    .pipe concat 'main.min.js'
    .pipe(if IS_PROD then minifier minifierOption else gutil.noop())
    .pipe gulp.dest ''
    .pipe reload {stream: true}

gulp.task 'browser-sync', ->
  browserSync
    ui: false
    proxy: 'localhost:3000'
    port: 3002

gulp.task 'watch', ['build-less', 'coffee', 'build-js', 'browser-sync'], ->  
  gulp.watch 'css/*.less', ['build-less']
  gulp.watch coffeeFiles, ['coffee']
  gulp.watch jsFiles, ['build-js']

gulp.task 'prod', ->
  IS_PROD = true
  runSequence 'build-less', 'coffee', 'build-js'

gulp.task 'default', ['build-less', 'coffee', 'build-js']

// Pre-requisites: need to install all the npm modules with:
// npm install

// TODO:
// - use gulp-uglify for prod to minifiy javascript:
//   var uglify= require('gulp-uglify');
//   .pipe(uglify())
// - concat js files see http://www.hongkiat.com/blog/getting-started-with-gulp-js/

var browserify  = require('browserify');
var exorcist    = require('exorcist')
var gulp        = require('gulp');
var prefix      = require('gulp-autoprefixer');
var concat      = require('gulp-concat'); // TODO: this is not being used maybe remove?
var htmlreplace = require('gulp-html-replace');
var uglify      = require('gulp-uglify');
var react       = require('gulp-react');
var rename      = require('gulp-rename');
var streamify   = require('gulp-streamify');
var source      = require('vinyl-source-stream');
var buffer      = require('vinyl-buffer')
var watchify    = require('watchify');
var babelify    = require("babelify");

gulp.task('js', function() {
  browserify({
    entries: ['jsx/App.jsx'],
    debug: true
  })
    .transform('babelify', {presets: ['es2015', 'react']})
    .bundle()
    .pipe(exorcist('s/dist/bundle.js.map'))
    .pipe(source('bundle.js'))
    .pipe(gulp.dest('s/dist'))
});

// TODO: doesn't work
// http://stackoverflow.com/questions/24992980/how-to-uglify-output-with-browserify-in-gulp
// https://wehavefaces.net/gulp-browserify-the-gulp-y-way-bb359b3f9623#.wbnxdlbi1
gulp.task('jsmin', function() {
  browserify({
    entries: ['jsx/App.jsx'],
    debug: true
  })
    .transform('babelify', {presets: ['es2015', 'react']})
    .bundle()
    //.pipe(exorcist('s/dist/bundle.min.js.map'))
    .pipe(source('bundle.min.js'))
    .pipe(buffer())
    .pipe(uglify())
    .pipe(gulp.dest('s/dist'))
});

gulp.task('css', function() {
  return gulp.src('s/*.css')
  .pipe(prefix('last 2 versions'))
  .pipe(gulp.dest('s/css/'));
});

gulp.task('watch', function() {
  gulp.watch(['jsx/*', 's/*.css'], ['css', 'js']);
});

gulp.task('build_and_watch', ['css', 'js', 'watch']);

gulp.task('default', ['css', 'js']);

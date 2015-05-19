// Pre-requisites: need to install all the npm modules with:
// npm install -g gulp-rename (etc.)
// If you don't have all the modules installed locally, run as:
// NODE_PATH=/usr/local/lib/node_modules gulp css

// TODO:
// - convert webpack.config.js to gulp,  http://tylermcginnis.com/reactjs-tutorial-pt-2-building-react-applications-with-gulp-and-browserify/
// - use gulp-uglify for prod to minifiy javascript:
//   var uglify= require('gulp-uglify');
//   .pipe(uglify())
// - concat js files see http://www.hongkiat.com/blog/getting-started-with-gulp-js/

var gulp = require('gulp');
var prefix = require('gulp-autoprefixer');
var rename = require('gulp-rename');
var concat = require('gulp-concat');
var uglify = require('gulp-uglify');
var react = require('gulp-react');
var htmlreplace = require('gulp-html-replace');
var source = require('vinyl-source-stream');
var browserify = require('browserify');
var watchify = require('watchify');
var reactify = require('reactify');
var streamify = require('gulp-streamify');
var exorcist   = require('exorcist')

gulp.task('js', function() {
  browserify({
    entries: ['jsx/App.jsx'],
    transform: [reactify],
    debug: true
  })
    .bundle()
    .pipe(exorcist('s/dist/bundle.js.map'))
    .pipe(source('bundle.js'))
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

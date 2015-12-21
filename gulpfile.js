// Pre-requisites: need to install all the npm modules with:
// npm install

var browserify  = require('browserify');
var exorcist    = require('exorcist');
var gulp        = require('gulp');
var prefix      = require('gulp-autoprefixer');
var uglify      = require('gulp-uglify');
var react       = require('gulp-react');
var sass        = require('gulp-sass');
var sourcemaps  = require('gulp-sourcemaps');
var source      = require('vinyl-source-stream');
var buffer      = require('vinyl-buffer')
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

gulp.task('jsmin', function() {
  browserify({
    entries: ['jsx/App.jsx'],
    debug: true
  })
    .transform('babelify', {presets: ['es2015', 'react']})
    .bundle()
    .pipe(exorcist('s/dist/bundle.min.js.map'))
    .pipe(source('bundle.min.js'))
    .pipe(buffer())
    .pipe(uglify())
    .pipe(gulp.dest('s/dist'))
});

gulp.task('css', function() {
  return gulp.src('./sass/main.scss')
  .pipe(sourcemaps.init())
  .pipe(sass().on('error', sass.logError))
  .pipe(prefix('last 2 versions'))
  .pipe(sourcemaps.write('.')) // this is relative to gulp.dest()
  .pipe(gulp.dest('./s/dist/'));
});

gulp.task('watch', function() {
  gulp.watch('jsx/**', ['js']);
  gulp.watch(['sass/**/*'], ['css']);
});

gulp.task('build_and_watch', ['css', 'js', 'watch']);

gulp.task('default', ['css', 'js']);

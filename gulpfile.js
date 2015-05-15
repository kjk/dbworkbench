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

gulp.task('css', function() {
  return gulp.src('s/*.css')
  .pipe(prefix('last 2 versions'))
  .pipe(gulp.dest('s/css/'));
});

gulp.task('watch', function() {
  gulp.watch('s/*.css', ['css']);
});

gulp.task('css_and_watch', ['css', 'watch']);

gulp.task('default', ['css']);

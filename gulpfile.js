var gulp = require('gulp');
var sass = require('gulp-sass');
var autoprefixer = require('gulp-autoprefixer');
var sourcemaps = require('gulp-sourcemaps');
var del = require('del');
var browserSync = require('browser-sync').create();

var http_port = 3000;
var s_dir = "./src";
var d_dir = "./dist";


var paths = {
  bootstrap_assets: './bower_components/bootstrap-sass/assets' ,
  jquery_assets: './bower_components/jquery/dist' ,

  sass: {
    watch: [s_dir + '/scss/**/*.{scss,sass}'],
    src: s_dir + '/scss/**/*.{scss,sass}',
    dest: d_dir + '/css'
  },

  html: {
    watch: [s_dir + '/html/**/*.html'],
    src: s_dir + '/html/**/*.html',
    dest: d_dir + '/html'
  },

  font: {
    dest: d_dir + '/font'  /* must match $icon-font-path variable in bootstrap.scss */
  },

  js: {
    dest: d_dir + '/js'
  },

  img: {
    src: s_dir + '/img',
    dest: d_dir + '/img'
  }
};

/**
 * CLEAN
 */

// Not all tasks need to use streams
// A gulpfile is just another node program and you can use any package available on npm
gulp.task('clean', function() {
  // You can use multiple globbing patterns as you would with `gulp.src`
  return del([d_dir]);
});


/*
cat
js/bootstrap-transition.js
js/bootstrap-alert.js
js/bootstrap-button.js
js/bootstrap-carousel.js
js/bootstrap-collapse.js
js/bootstrap-dropdown.js
js/bootstrap-modal.js
js/bootstrap-tooltip.js
js/bootstrap-popover.js
js/bootstrap-scrollspy.js
js/bootstrap-tab.js
js/bootstrap-typeahead.js
> bootstrap/js/bootstrap.js
*/



/**
 * JS
 */
gulp.task('js-bootstrap', function() {
  return gulp.src([
      paths.bootstrap_assets + '/javascripts/bootstrap.js',
      paths.bootstrap_assets + '/javascripts/bootstrap.min.js'
    ])
    .pipe(gulp.dest(paths.js.dest));
});

gulp.task('js-jquery', function() {
  return gulp.src([
      paths.jquery_assets + '/jquery.js',
      paths.jquery_assets + '/jquery.min.js'
    ])
    .pipe(gulp.dest(paths.js.dest));
});

gulp.task('js', ['js-jquery', 'js-bootstrap']);

/**
 * FONT
 */
gulp.task('font-bootstrap', function() {
  return gulp.src(paths.bootstrap_assets + '/fonts/**/*.*')
    .pipe(gulp.dest(paths.font.dest));
});

gulp.task('font', ['font-bootstrap']);


/**
 * IMG
 */
gulp.task('img', function() {
  gulp.src(paths.img.src + "/favicon.ico")
    .pipe(gulp.dest(d_dir));
  gulp.src([paths.img.src + "/**/*", "!**/favicon.ico"])
    .pipe(gulp.dest(paths.img.dest));
});


/**
 * STATIC
 */
 gulp.task('static', ['font', 'js', 'img']);


/**
 * SASS
 * With includePaths option our main style.scss can import bootstrap easily. Here is an example:
 *    @import "bootstrap";
 */

gulp.task('sass', function() {
  return gulp.src(paths.sass.src)
    .pipe(sourcemaps.init({ loadMaps: true }))
    .pipe(sass({
      includePaths: [
        paths.sass.src,
        paths.bootstrap_assets + '/stylesheets'
      ]
    }).on('error', sass.logError))
    .pipe(autoprefixer())
    .pipe(sourcemaps.write('./'))
    .pipe(gulp.dest(paths.sass.dest))
    .pipe(browserSync.stream({match: '**/*.css'}));
});


/**
 * HTML
 */

gulp.task('html', function() {
  return gulp.src(paths.html.src)
    // Perform minification tasks, etc here
    .pipe(gulp.dest(paths.html.dest))
    .pipe(browserSync.stream({match: '**/*.html'}));
});


/**
 * SERVE
 */

// Static Server + watching scss/html files
gulp.task('serve', ['sass', 'html'], function() {

    browserSync.init({
      // Serve files from the app directory, with a specific index filename
      server: {
          baseDir: [d_dir + "/html", d_dir],
          index: "index.html"
      },
      port: http_port
    });

    gulp.watch(paths.sass.watch, ['sass']);
//    gulp.watch(paths.html.watch).on('change', browserSync.reload);
    gulp.watch(paths.html.watch, ['html']);
});

/**
 * DEFAULT
 */

 gulp.task('default',  ['serve']);

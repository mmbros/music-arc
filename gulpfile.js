// http://struct.cc/blog/2015/05/08/building-web-applications-in-golang-with-gulp-and-livereload/
// https://gist.github.com/squidfunk/120b6f02927fdc9ef9f1

var gulp = require('gulp');
var sass = require('gulp-sass');
var child = require('child_process');
var autoprefixer = require('gulp-autoprefixer');
var sourcemaps = require('gulp-sourcemaps');
var del = require('del');
var browserSync = require('browser-sync').create();
var util = require('gulp-util');
//var notifier   = require('node-notifier');
var sync = require('gulp-sync')(gulp).sync;

//var reload  = require('gulp-livereload');
var reload      = browserSync.reload;


var http_port = 3000;
var s_dir = "./src";
var d_dir = "./dist";


var paths = {
  bootstrap: './bower_components/bootstrap-sass/assets' ,
  jquery: './bower_components/jquery/dist' ,

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
      paths.bootstrap + '/javascripts/bootstrap.js',
      paths.bootstrap + '/javascripts/bootstrap.min.js'
    ])
    .pipe(gulp.dest(paths.js.dest));
});

gulp.task('js-jquery', function() {
  return gulp.src([
      paths.jquery + '/jquery.js',
      paths.jquery + '/jquery.min.js'
    ])
    .pipe(gulp.dest(paths.js.dest));
});

gulp.task('js', ['js-jquery', 'js-bootstrap']);

/**
 * FONT
 */
gulp.task('font-bootstrap', function() {
  return gulp.src(paths.bootstrap + '/fonts/**/*.*')
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


 /* ----------------------------------------------------------------------------
  * Assets pipeline
  * ------------------------------------------------------------------------- */

/**
 * SASS
 * With includePaths option our main style.scss can import bootstrap easily. Here is an example:
 *    @import "bootstrap";
 */

gulp.task('assets:stylesheets', function() {
  return gulp.src(paths.sass.src)
    .pipe(sourcemaps.init({ loadMaps: true }))
    .pipe(sass({
      includePaths: [
        paths.sass.src,
        paths.bootstrap + '/stylesheets'
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

/*
 * Build assets.
 */
gulp.task('assets:build', [
  'assets:stylesheets'
]);

/*
 * Watch assets for changes and rebuild on the fly.
 */
gulp.task('assets:watch', function() {

  /* Rebuild stylesheets on-the-fly */
  gulp.watch([
    'src/scss/**/*.{sass,scss}'
  ], ['assets:stylesheets']);

});

 /* ----------------------------------------------------------------------------
  * GO
  * ------------------------------------------------------------------------- */

 /* GO Application server */
 var server = null;

/*
 * Build application server: go install
 */
gulp.task('server:build', function() {
  var build = child.spawnSync('go', ['install']);
  if (build.stderr.length) {
    var lines = build.stderr.toString()
      .split('\n').filter(function(line) {
        return line.length
      });
    for (var l in lines)
      util.log(util.colors.red(
        'Error (go install): ' + lines[l]
      ));
/*
    notifier.notify({
      title: 'Error (go install)',
      message: lines
    });
*/
  }
  return build;
});

/*
 * Create source files: go generate
 */
gulp.task('server:generate', function() {
  var build = child.spawnSync('go', ['generate']);
  if (build.stderr.length) {
    var lines = build.stderr.toString()
      .split('\n').filter(function(line) {
        return line.length
      });
    for (var l in lines)
      util.log(util.colors.red(
        'Error (go generate): ' + lines[l]
      ));
  }
  return build;
});

/*
 * Restart application server.
 */
gulp.task('server:spawn', function() {
  if (server){
    server.kill();
  }

  /* Spawn application server */
  server = child.spawn('music-arc');

  /* Trigger reload upon server start */
/*
  server.stdout.once('data', function() {
    // reload.reload('/');
    reload();
  });
*/
  /* Pretty print server log output */
  server.stdout.on('data', function(data) {
    var lines = data.toString().split('\n')
    for (var l in lines)
      if (lines[l].length)
        util.log(lines[l]);
  });

  /* Print errors to stdout */
  server.stderr.on('data', function(data) {
    process.stdout.write(data.toString());
  });
});


/*
 * Watch source for changes and restart application server.
 */
gulp.task('server:watch', function() {

  /* Restart application server */
  gulp.watch([
    'templates/**/*.tmpl',
    'data/music-arc-inc.xml'
  ], ['server:spawn']);

  /* Rebuild and restart application server */
  gulp.watch([
    '**/*.go',
  ], sync([
    'server:build',
    'server:spawn'
  ], 'server'));

  /* Re Generate source files */
  gulp.watch([
    'templates.config.toml',
  ], ['server:generate']);
});



/* ----------------------------------------------------------------------------
 * Interface
 * ------------------------------------------------------------------------- */

/*
 * Build assets and application server.
 */
gulp.task('build', [
  'assets:build',
  'server:generate',
  'server:build'
]);


/*
 * Start asset and server watchdogs and initialize livereload.
 */
gulp.task('watch', [
  'assets:build',
  'server:build'
], function() {

//  reload.listen();

  browserSync.init();


  return gulp.start([
    'assets:watch',
    'server:watch',
    'server:spawn'
  ]);
});

/*
 * Build assets by default.
 */
gulp.task('default', ['build']);

/**
 * SERVE
 */

// Static Server + watching scss/html files
gulp.task('serve__OLD', ['sass', 'html'], function() {

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

# Web project template

Create a web project with `nodejs`, `gulp`, `git`, `sass`.

1. Create project folder and change directory.
   The `$_` (dollar underscore) bash command variable contains the most recent parameter.
       mkdir <appdir> && cd $_

2. Create and initialize `package.json`
       npm init

3. Install `gulp` and save in `devDependencies` section of `package.json`
       npm install gulp -D
   Install some packages for gulp tasks
       npm install browser-sync gulp-sass gulp-autoprefixer gulp-sourcemaps del -D
   Don't mind if eventually will occur the warning:
       npm WARN prefer global node-gyp@3.2.0 should be installed with -g

4. `bower` configuration for Bootstrap
      # npm install bower -D
      bower init
      bower install bootstrap-sass -D


5. Create the `gulpfile.js`
       touch gulpfile.js

6. Edit `gulpfile.js` as:
       var gulp = require('gulp')
       var sass = require('gulp-sass')
       var autoprefixer = require('gulp-autoprefixer')
       var sourcemaps = require('gulp-sourcemaps')

6. Inizialize `git` repository and create `.gitignore` file
       git init
       echo "/node_modules" > .gitignore

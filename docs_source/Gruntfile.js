module.exports = function(grunt) {

    // 1. All configuration goes here
    grunt.initConfig({
        pkg: grunt.file.readJSON('package.json'),

        concat: {
            js: {
                src: ['js/**/*.js'],
                dest: 'prod/js/script.js'
            },
            css: {
                src: ['prod/css/style.css', 'prod/css/*', 'css_depend/**/*.css'], // The middle line is for ordering.. apparently.
                dest: 'prod/css/style.css'
            }
        },

        less: {
          production: {
            options: {
              cleancss: true,
              sourceMap: true
            },
            files: {
              "prod/css/style.css": "less/style.less"
            }
          }
        },

        uglify: {
            build: {
                src: 'prod/js/script.js',
                dest: 'prod/js/script.min.js'
            }
        },

        imagemin: {
            dynamic: {
                files: [{
                    expand: true,
                    cwd: 'img/',
                    src: ['**/*.{png,jpg,gif}'],
                    dest: 'prod/img'
                }]
            }
        },

        htmlmin: {                                     // Task
            dist: {                                      // Target
                options: {                                 // Target options
                    removeComments: true,
                    collapseWhitespace: true,
                    removeCommentsFromCDATA: true,
                    removeCDATASectionsFromCDATA: true,
                    collapseBooleanAttributes: true,
                    removeAttributeQuotes: true,
                    removeRedundantAttributes: true,
                    useShortDoctype: true,
                    removeEmptyAttributes: true,

              },
              files: {
                    'prod/index.html': 'index.html',     // 'destination': 'source'
                    'prod/how.html' : 'how.html',
                    'prod/about.html' : 'about.html',
                    'prod/download.html' : 'download.html',
                    'prod/blog.html' : 'blog.html',
                    'prod/license.html' : 'license.html',
                    'prod/legal/license.html' : 'legal/license.html',
                    'prod/update.html' : 'update.html',
                    'prod/updatecheck.html' : 'updatecheck.html',
                    'prod/sending_logs.html' : 'sending_logs.html',
                    'prod/linux_download.html' : 'linux_download.html',

                    // this is thorny as fuck to fix for global case.
              }
            }
        },

        cssmin: {
          minify: {
            expand: true,
            src: ['prod/css/style.css'],
          }
        },

        watch: {
            options: {
                    livereload: true,
                },
            scripts: {
                files: ['js/*.js'],
                tasks: ['concat:js', 'uglify'],
                options: {
                    spawn: false,
                },
            },
            css: {
                files: ['less/*.less'],
                tasks: ['less', 'autoprefixer', 'concat:css', 'cssmin'],
                options: {
                    spawn: false,
                }
            },
            html: {
                files: ['*.html'],
                tasks: ['htmlmin'],
                options: {
                    spawn: false,
                }
            },
        },

        autoprefixer: {
            options: {
              // Task-specific options go here.
            },
            single_file: {
                options: {
                    // Target-specific options go here.
                },
                src: 'prod/css/style.css',
                dest: 'prod/css/style.css'
                },
        },

    });

    // 3. Where we tell Grunt we plan to use this plug-in.
    require('load-grunt-tasks')(grunt);

    // 4. Where we tell Grunt what to do when we type "grunt" into the terminal.
    grunt.registerTask('default', ['less', 'uglify', 'htmlmin', 'concat', 'cssmin', 'autoprefixer', 'watch']);

};
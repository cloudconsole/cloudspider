module.exports = function (grunt) {

    grunt.initConfig({
        pkg: grunt.file.readJSON('package.json'),

        clean: ["ui/static/src/build"],

        copy: {
            main: {
                src: 'ui/static/src/build/cloudspider.min.js',
                dest: 'ui/static/js/cloudspider.js',
            },
        },

        watch: {
            files: ['ui/static/src/*.js'],
            tasks: ['clean', 'browserify', 'uglify', 'copy', 'clean']
        },

        browserify: {
            dist: {
                options: {
                    transform: [['babelify', {presets: ['es2015', 'react']}]]
                },
                // the files to concatenate
                src: ['ui/static/src/main.js'],
                // the location of the resulting JS file
                dest: 'ui/static/src/build/bundle.js',
            }
        },

        uglify: {
            options: {
                // the banner is inserted at the top of the output
                banner: '/*! <%= pkg.name %> <%= grunt.template.today("dd-mm-yyyy") %> */\n'
            },
            dist: {
                files: {
                    'ui/static/src/build/<%= pkg.name %>.min.js': ['<%= browserify.dist.dest %>']
                }
            }
        }

    });

    grunt.loadNpmTasks('grunt-contrib-clean');
    grunt.loadNpmTasks('grunt-contrib-copy');
    grunt.loadNpmTasks('grunt-contrib-uglify');
    grunt.loadNpmTasks('grunt-contrib-cssmin');
    grunt.loadNpmTasks('grunt-contrib-watch');
    grunt.loadNpmTasks('grunt-browserify');

    // this would be run by typing "grunt test" on the command line
    // grunt.registerTask('test', ['watch']);

    // the default task can be run just by typing "grunt" on the command line
    grunt.registerTask('default', ['clean', 'browserify', 'uglify', 'copy', 'clean']);

};

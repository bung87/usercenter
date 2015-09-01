/*global module:false*/
module.exports = function(grunt) {

  // Project configuration.
  grunt.initConfig({
    // Metadata.
    pkg: grunt.file.readJSON('package.json'),
    banner: '/*! <%= pkg.title || pkg.name %> - v<%= pkg.version %> - ' +
    '<%= grunt.template.today("yyyy-mm-dd") %>\n' +
    '<%= pkg.homepage ? "* " + pkg.homepage + "\\n" : "" %>' +
    '* Copyright (c) <%= grunt.template.today("yyyy") %> <%= pkg.author.name %>;' +
    ' Licensed <%= _.pluck(pkg.licenses, "type").join(", ") %> */\n',
    // Task configuration.
    concat: {
      options: {
        banner: '<%= banner %>',
        stripBanners: true
      },
      dist: {
        src: ['src/js/<%= pkg.name %>.js'],
        dest: 'public/js/<%= pkg.name %>.js'
      }
    },
    uglify: {
      options: {
        banner: '<%= banner %>'
      },
      dist: {
        src: '<%= concat.dist.dest %>',
        dest: 'public/js/<%= pkg.name %>.min.js'
      }
    },
/*    watch: {
    
      lib_test: {
        files: '<%= jshint.lib_test.src %>'
      }
    },*/
    bower_concat: {
      all: {
        dest: 'public/js/venders.pkgd.js',
        cssDest: 'public/css/venders.css',
        exclude: [
        'jquery'
        ],
        bowerOptions: {
          relative: false
        }
      }
    },
    bowercopy: {
      options: {
        // Bower components folder will be removed afterwards 
        clean: false
      },
    // Anything can be copied 

    // Javascript 
    libs: {
      options: {
        destPrefix: 'public/js'
      },
      files: {
        'jquery.min.js': 'jquery/dist/jquery.min.js'
            // 'jquery.min.map': 'jquery/dist/jquery.min.map'
            // 'zepto.min.js':'zepto/zepto.min.js'
          }

        },
        plugins: {
          options: {
            destPrefix: 'public/plugins'
          }
        },

    // Images 
    images: {
      options: {
        destPrefix: 'public/images'
      }
    },

  }
});

  // These plugins provide necessary tasks.
  grunt.loadNpmTasks('grunt-contrib-concat');
  grunt.loadNpmTasks('grunt-contrib-uglify');
 
  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-bower-concat');
  grunt.loadNpmTasks('grunt-bowercopy');
  // Default task.
  grunt.registerTask('default', ['bower_concat:all','bowercopy:libs', 'concat', 'uglify']);

};

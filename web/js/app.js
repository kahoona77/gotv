'use strict';

angular.module('xtv.controllers', []);
angular.module('xtv.services', []);

// Declare app level module which depends on filters, and services
angular.module('xtv', [
  'ngRoute',
  'ngAnimate',
  'xtv.filters',
  'xtv.services',
  'xtv.directives',
  'xtv.controllers',
]).config(['$routeProvider', '$locationProvider', function($routeProvider, $locationProvider) {

  $locationProvider.html5Mode(true);

  $routeProvider.when('/home', {templateUrl: 'web/partials/home.html', controller: 'HomeCtrl'});
  $routeProvider.when('/search', {templateUrl: 'web/partials/search.html', controller: 'SearchCtrl'});
  $routeProvider.when('/downloads', {templateUrl: 'web/partials/downloads.html', controller: 'DownloadsCtrl'});
  $routeProvider.when('/logFile', {templateUrl: 'web/partials/logFile.html', controller: 'LogFileCtrl'});
  $routeProvider.otherwise({redirectTo: '/home'});
}]).run(['$rootScope', function ($rootScope) {
        $rootScope.showSettingsDialog = function () {
            $rootScope.$broadcast ('xtv:showSettingsDialog');
        }
    }]);

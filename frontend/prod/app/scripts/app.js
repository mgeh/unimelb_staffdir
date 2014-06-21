'use strict';

angular
  .module('staffdirApp', [
    'ngCookies',
    'ngResource',
    'ngSanitize',
    'ngRoute'
  ])
  .config(function ($routeProvider) { // $httpProvider
    // delete $httpProvider.defaults.headers.common['X-Requested-With'];
    // $httpProvider.defaults.headers.common['Access-Control-Allow-Origin'] =  '*';
    $routeProvider
      .when('/search/:query?', {
        templateUrl: 'views/main.html',
        controller: 'MainCtrl'
      })
      .when('/person/:personId', {
        templateUrl: 'views/person.html',
        controller: 'PersonCtrl'
      })
      .when('/', {
        templateUrl: 'views/main.html',
        controller: 'MainCtrl'
      })
      .otherwise({
        redirectTo: '/search'
      });
  });


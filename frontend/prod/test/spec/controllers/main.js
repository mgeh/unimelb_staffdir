'use strict';

describe('Controller: MainCtrl', function () {

  // load the controller's module
  beforeEach(module('staffdirApp'));

  var MainCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    MainCtrl = $controller('MainCtrl', {
      $scope: scope
    });
  }));

  it('should have an empty data array upon init', function () {
    expect(scope.data.length).toBe(0);
  });
});

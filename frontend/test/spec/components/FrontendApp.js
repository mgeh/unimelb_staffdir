'use strict';

describe('Main', function () {
  var FrontendApp, component;

  beforeEach(function () {
    var container = document.createElement('div');
    container.id = 'content';
    document.body.appendChild(container);

    FrontendApp = require('../../../src/scripts/components/FrontendApp.jsx');
    component = FrontendApp();
  });

  it('should create a new instance of FrontendApp', function () {
    expect(component).toBeDefined();
  });
});

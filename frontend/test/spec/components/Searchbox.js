'use strict';

describe('Searchbox', function () {
  var Searchbox, component;

  beforeEach(function () {
    Searchbox = require('../../../src/scripts/components/Searchbox.jsx');
    component = Searchbox();
  });

  it('should create a new instance of Searchbox', function () {
    expect(component).toBeDefined();
  });
});

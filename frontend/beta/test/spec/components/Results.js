'use strict';

describe('Results', function () {
  var Results, component;

  beforeEach(function () {
    Results = require('../../../src/scripts/components/Results.jsx');
    component = Results();
  });

  it('should create a new instance of Results', function () {
    expect(component).toBeDefined();
  });
});

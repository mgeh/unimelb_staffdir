'use strict';

describe('PersonDetails', function () {
  var PersonDetails, component;

  beforeEach(function () {
    PersonDetails = require('../../../src/scripts/components/PersonDetails.jsx');
    component = PersonDetails();
  });

  it('should create a new instance of PersonDetails', function () {
    expect(component).toBeDefined();
  });
});

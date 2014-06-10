'use strict';

describe('Person', function () {
  var Person, component;

  beforeEach(function () {
    Person = require('../../../src/scripts/components/Person.jsx');
    component = Person();
  });

  it('should create a new instance of Person', function () {
    expect(component).toBeDefined();
  });
});

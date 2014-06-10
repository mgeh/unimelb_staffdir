'use strict';

describe('PersonList', function () {
  var PersonList, component;

  beforeEach(function () {
    PersonList = require('../../../src/scripts/components/PersonList.jsx');
    component = PersonList();
  });

  it('should create a new instance of PersonList', function () {
    expect(component).toBeDefined();
  });
});

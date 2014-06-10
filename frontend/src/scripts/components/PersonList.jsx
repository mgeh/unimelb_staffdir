/**
 * @jsx React.DOM
 */

'use strict';

var React = require('react/addons');
require('../../styles/PersonList.css');
var Person = require("./Person.jsx");

var PersonList = React.createClass({
  render: function () {
  	var people = this.props.results.map(function(person) {
  		return <Person name={person.name} />
  	});
    return (
        <div>
          {people}
        </div>
      );
  }
});

module.exports = PersonList;

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
      // console.log(person);
  		return <Person key={person['email']} person={person} />
  	});
    return (
        <div>
          {people}
        </div>
      );
  }
});

module.exports = PersonList;

/**
 * @jsx React.DOM
 */

'use strict';

var React = require('react/addons');
require('../../styles/Person.css');
var Details = require("./PersonDetails.jsx");


var Person = React.createClass({
	getDetails: function(email) {
		return (
			<div>{email}</div>
			);
	},

  render: function () {
  	var scope = this;
  	function clickOn(email){
  		console.log();
  		React.renderComponent(
			  React.DOM.div(
			    null,
			    <Details person={scope.props.person}/>
			  ),
			  document.body
			);
  	};
    return (
        <div>
          <p><a href="#" onClick={clickOn}>{this.props.person["a.name"]}</a>, {this.props.person["a.position"]}, {this.props.person["a.department"]}, {this.props.person["a.email"]}, {this.props.person["a.phone"]}</p>
        </div>
      );
  }
});

module.exports = Person;

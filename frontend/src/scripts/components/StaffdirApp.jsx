/**
 * @jsx React.DOM
 */

'use strict';

var React = require('react/addons');
var ReactTransitionGroup = React.addons.TransitionGroup;

// CSS
require('../../styles/reset.css');
require('../../styles/main.css');

// All compontents
var Searchbox = require("./Searchbox.jsx");
var Results = require("./Results.jsx");

var imageURL = '../../images/yeoman.png';

var results = [
  {name: 'Val Lyashov', dep: 'ITS'},
  {name: 'Tom Stringer', dep: 'ITS'}
];


var StaffdirApp = React.createClass({

	getData: function(query){
		$.ajax({
		      url: 'http://uom-13melb.herokuapp.com/area/159',
		      dataType: 'json',
		      data: query,
		      success: function(data) {
		        this.setState({data: query});
		      }.bind(this),
		      error: function(xhr, status, err) {
		        console.error(this.props.url, status, err.toString());
		      }.bind(this)
		})
	},
  render: function() {
    return (
      <div className='main'>
        <Searchbox />
        <Results results={this.props.results} />
      </div>

    );
  }
});

React.renderComponent(<StaffdirApp results={results}/>, document.getElementById('content')); // jshint ignore:line
module.exports = StaffdirApp;

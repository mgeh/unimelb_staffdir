/**
 * @jsx React.DOM
 */

'use strict';

var React = require('react/addons');
// var Router = require('react-router-component')

var ReactTransitionGroup = React.addons.TransitionGroup;

// CSS
require('../../styles/reset.css');
require('../../styles/main.css');

// All compontents
var Searchbox = require("./Searchbox.jsx");
var Results = require("./Results.jsx");

var imageURL = '../../images/yeoman.png';
var url_base = 'http://uom-staffdir.herokuapp.com/staffdir/person/';

var StaffdirApp = React.createClass({

	handleSubmit: function(e){
		// this..preventDefault();
		$.ajax({
			url: url_base + e,
			dataType: "json",
			success: function (data) {
				this.setState({results: data});
				// console.log(this.state.results);
			}.bind(this)
		});
	},
	getInitialState: function(){
		return {results:[]}
	},
  render: function() {
  	
    return (
      <div className='main'>
        <Searchbox results={this.state.results} getData={this.handleSubmit} />
        <Results results={this.state.results}  />
      </div>

    );
  }
});

React.renderComponent(<StaffdirApp />, document.getElementById('content')); // jshint ignore:line
module.exports = StaffdirApp;

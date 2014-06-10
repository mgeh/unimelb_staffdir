/**
 * @jsx React.DOM
 */

'use strict';

var React = require('react/addons');
require('../../styles/Searchbox.css');

var url_base = 'http://uom-staffdir-api.elasticbeanstalk.com/staffdir/person/';
var Searchbox = React.createClass({
	getInitialState: function(){
		return {query:''}
	},
	handleSubmit: function(e){
		e.preventDefault();
		$.ajax({
			url: url_base,
			dataType: "json",
			success: function (data) {
				this.setState({results: data});
			}.bind(this)
		});
	},
	searchChange: function(e){
		this.setState({query:e.target.value});
	},
  render: function () {
    return (
        <div>
          <h2>Search </h2>
          <form onSubmit={this.handleSubmit}>
          	<input type='text' value={this.state.query} onChange={this.searchChange} />
          	<button>search</button>
          </form>
        </div>
      );
  }
});

module.exports = Searchbox;

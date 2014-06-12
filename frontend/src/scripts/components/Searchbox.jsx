/**
 * @jsx React.DOM
 */

'use strict';

var React = require('react/addons');
require('../../styles/Searchbox.css');


var Searchbox = React.createClass({
	getInitialState: function(){
		return {query:''}
	},

	searchChange: function(e){
		this.setState({query:e.target.value});
	},
	handleSubmit: function() {
	    var query = this.refs.query.getDOMNode().value.trim() || '';
	    console.log(query);
	    this.props.getData(query);
	    this.refs.query.getDOMNode().value = '';
	    return false;
	},
  render: function () {
    return (
        <div>
          <h2>Search </h2>
          <form onSubmit={this.handleSubmit}>
          	<input type='text' ref="query" onChange={this.searchChange} />
          	<button>search</button>
          </form>
        </div>
      );
  }
});

module.exports = Searchbox;

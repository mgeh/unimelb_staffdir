/**
 * @jsx React.DOM
 */

'use strict';

var React = require('react/addons');
require('../../styles/Results.css');
var PersonList = require("./PersonList.jsx");


var Results = React.createClass({
  render: function () {
    return (
        <div>
          <p>Results</p>
          <PersonList results={this.props.results} />
        </div>
      );
  }
});

module.exports = Results;

/**
 * @jsx React.DOM
 */

'use strict';

var React = require('react/addons');
require('../../styles/Person.css');

var Person = React.createClass({

  render: function () {
    return (
        <div>
          <p>{this.props.name}</p>
        </div>
      );
  }
});

module.exports = Person;

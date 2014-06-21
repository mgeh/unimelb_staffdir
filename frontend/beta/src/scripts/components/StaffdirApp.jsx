/**
 * @jsx React.DOM
 */

'use strict';

var React = require('react/addons');
var Router = require('react-router-component');

var Locations = Router.Locations;
var Location = Router.Location;
var NotFound = Router.NotFound;
var Link = require('react-router-component').Link;

var ReactTransitionGroup = React.addons.TransitionGroup;

// CSS
require('../../styles/reset.css');
require('../../styles/main.css');

// All compontents
var Searchbox = require("./Searchbox.jsx");
var Results = require("./Results.jsx");

var imageURL = '../../images/yeoman.png';
var url_base = 'http://uom-staffdir.herokuapp.com/staffdir/person/';

var MainPage = React.createClass({
  getInitialState: function(){
    return {results:[]}
  },
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
  
  render: function() {
      return (
        <div className='main'>
          <Searchbox results={this.state.results} getData={this.handleSubmit} />
          <Results results={this.state.results}  />
        </div>
        );
    }
  });

var NotFoundPage = React.createClass({
  render: function(){
    return (
      <div>
        <p>Error occured with pathing</p>
        <Link href="/">user page</Link>.
        <Link href="/about">about</Link>.
      </div>
    )
  }
});

var AboutPage = React.createClass({
  render: function(){
    return (
      <div>
        <p>Error occured with pathing</p>
        <Link href="/">user page</Link>.
      </div>
    )
  }
});

var StaffdirApp = React.createClass({

  render: function() {
    return (
      <Locations>
        <Location path="/" handler={MainPage} />
        <Location path="/about" handler={AboutPage} />
        <NotFound handler={NotFoundPage} />
      </Locations>

    );
  }
});

React.renderComponent(<StaffdirApp />, document.getElementById('content')); // jshint ignore:line
module.exports = StaffdirApp;

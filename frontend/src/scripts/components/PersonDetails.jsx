/**
 * @jsx React.DOM
 */

'use strict';

var React = require('react/addons');
require('../../styles/PersonDetails.css');
var Person = require("./Person.jsx");

var PersonDetails = React.createClass({
	getInitialState: function() {
		return ({colleagues: [],
						 reports: [],
						 manager: []});
	},

	getColleagues: function(){
		$.ajax({
			url: 'http://uom-staffdir.herokuapp.com/staffdir/colleagues/' + this.props.person["a.email"],
			dataType: "json",
			success: function (data) {
				this.setState({colleagues: data});
			}.bind(this)
		});
	},

	getReports: function(){
		$.ajax({
			url: 'http://uom-staffdir.herokuapp.com/staffdir/reports/' + this.props.person["a.email"],
			dataType: "json",
			success: function (data) {
				this.setState({reports: data});
			}.bind(this)
		});
	},

	getManager: function(){
		$.ajax({
			url: 'http://uom-staffdir.herokuapp.com/staffdir/manager/' + this.props.person["a.email"],
			dataType: "json",
			success: function (data) {
				this.setState({manager: data});
			}.bind(this)
		});
	},

	componentWillMount : function() {
		this.getColleagues();
	  	this.getReports();
	 	this.getManager();
	},

  render: function () {
  	var scope = this;
  	
  	var colleagues = this.state.colleagues.map(function(colleague) {
      	console.log(colleague);
  		return <div>{colleague} </div>;
  	});

  	var managers = this.state.manager.map(function(manager) {
  		console.log(manager);
  		return (<div>{manager}</div>);
  	});

  	var reports = this.state.reports.map(function(report) {
      console.log(report);
  		return <div>{report} </div>;
  	});


    return (
			    <div>
			    	<h2>{this.props.person["a.name"]}</h2>
			    	<p>Email: {this.props.person["a.email"]}</p>
			    	<p>Phone: {this.props.person["a.phone"]}</p>
			    	<p>Manager: {managers}</p>
			    	<p>Colleagues:  {colleagues}</p>
			    	<p>Reports: {reports}</p>
			    </div>
      );
  }
});

module.exports = PersonDetails;

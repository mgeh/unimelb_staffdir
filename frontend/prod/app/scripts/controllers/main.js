'use strict';

angular.module('staffdirApp')
	.controller('MainCtrl', function ($scope, $http, $routeParams, $templateCache, $location) {
		$scope.method = 'GET';
		$scope.url = '';
		$scope.base = 'http://uom-staffdir.herokuapp.com/staffdir/person?q=';
		$scope.loading = false;
		$scope.message = '';
		$scope.data = [];
		$scope.params = $routeParams;
		$scope.suggestions = [];

		$scope.getTotal = function(dataset) {
			return dataset.length;
		};

		$scope.jumpTo = function(loc) {
			this['#c-'+loc.split('@')[0].split('.').join('-')].collapse('toggle');
		};

		$scope.fetch = function() {
			$scope.code = null;
			$scope.loading = true;
			$scope.response = null;
			$scope.message = '';
			$scope.data = [];
			$scope.suggestions = [];

			$http({method: $scope.method, url: $scope.base + $scope.url, cache: $templateCache}).
				success(function(data, status) {
					$location.path('/search/' + $scope.url);
					$scope.status = status;
					if (data.size === 100) {
						$scope.message = 'Reached search limit, returning first 100 results';
					} else if (data.size === 0) {
						$scope.message = 'Sorry, no results found for your query.';
					}
					$scope.loading = false;
					$scope.data = data;

					$scope.d2 = {};
					for(var i=0;i<$scope.data.data; i++) {
						if (!($scope.data.data[i].hostname in $scope.d2)) {
							$scope.d2[$scope.data.data[i].hostname] = {'data': [], 'count': 0};
						}
						$scope.d2[$scope.data.data[i].hostname].data.push($scope.data.data[i]);
						$scope.d2[$scope.data.data[i].hostname].count += 1;
					}
				}).
				error(function(data, status) {
					$scope.data = data || 'Request failed';
					$scope.message = 'Sorry, failed to retrieve the data.';
					$scope.status = status;
					$scope.loading = false;
				});
		};

		$scope.updateModel = function(method, url) {
			$scope.method = method;
			$scope.url = $scope.base + url;
		};

		$scope.tocsv = function (objArray) {
			var array = typeof objArray !== 'object' ? JSON.parse(objArray) : objArray;

			var str = '';
			var line = '';
			var j;
			if (this['#labels'].is(':checked')) {
				// var head = array[0];
				if (this['#quote'].is(':checked')) {
					for (j in array[0]) {
						var value = j + '';
						line += '"' + value.replace(/"/g, '""') + '",';
					}
				} else {
					for (j in array[0]) {
						line += j + ',';
					}
				}

				line = line.slice(0, -1);
				str += line + '\r\n';
			}

			for (var i = 0; i < array.length; i++) {
				line = '';

				if (this['#quote'].is(':checked')) {
					j = 0;
					for (j in array[i]) {
						var val = array[i][j] + '';
						line += '"' + val.replace(/"/g, '""') + '",';
					}
				} else {
					j = 0;
					for (j in array[i]) {
						line += array[i][j] + ',';
					}
				}

				line = line.slice(0, -1);
				str += line + '\r\n';
			}
			return str;
			
		};
				
			
		$scope.downloadcsv = function () {
			var csv = $scope.tocsv($scope.data.data);
			window.open('data:text/csv;charset=utf-8,' + this.escape(csv));
		};
		if ('query' in $scope.params) {
			$scope.url = $scope.params.query;
			$scope.fetch();
		}
	});

angular.module('staffdirApp')
	.controller('PersonCtrl', function ($scope, $http, $routeParams, $templateCache, $location) {
		$scope.method = 'GET';
		$scope.url = '';
		$scope.params = $routeParams;
		$scope.base = 'http://uom-staffdir.herokuapp.com/staffdir/details?id=' + $scope.params.personId;
		$scope.loading = false;
		$scope.message = '';
		$scope.data = [];
		$scope.suggestions = [];
		$scope.query = '';
		$scope.getTotal = function(dataset) {
			return dataset.length;
		};

		$scope.fetch = function() {
			$scope.code = null;
			$scope.loading = true;
			$scope.response = null;
			$scope.message = '';
			$scope.data = [];
			$scope.suggestions = [];

			$http({method: $scope.method, url: $scope.base, cache: $templateCache}).
				success(function(data, status) {
					$scope.status = status;
					if (data.size >= 100) {
						$scope.message = 'Reached search limit, returning first 100 results';
					} else if (data.size === 0) {
						return;
					}
					$scope.loading = false;
					$scope.data = data.data;

					$scope.d2 = {};
					for(var i=0;i<$scope.data.data; i++) {
						if (!($scope.data.data[i].hostname in $scope.d2)) {
							$scope.d2[$scope.data.data[i].hostname] = {'data': [], 'count': 0};
						}
						$scope.d2[$scope.data.data[i].hostname].data.push($scope.data.data[i]);
						$scope.d2[$scope.data.data[i].hostname].count += 1;
					}
				}).
				error(function(data, status) {
					$scope.data = data || 'Request failed';
					$scope.message = 'Sorry, failed to retrieve the data.';
					$scope.status = status;
					$scope.loading = false;
				});
		};
		$scope.fetch();

		$scope.search = function() {
			$location.path('/search/' + $scope.query);
		};

		$scope.goback = function() {
			window.history.back();
		};
	});


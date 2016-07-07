angular.module('simpleAPI', []);  

function mainController($scope, $http) {  
    //$scope.formData = {};

    // Cuando se cargue la página, pide del API todos los TODOs
    $http.get('http://127.0.0.1:8080/clinicaldata/validation')
        .success(function(data) {
            $scope.clinicalvalidationdataset = data;
			console.log(data)
        })
        .error(function(data, status) {
                console.log('Error status:' + status);
        });
		
	$http.get('http://192.168.248.128:8080/clinicaldata/training')
        .success(function(data) {
            $scope.clinicaltrainingdataset = data;
			console.log(data)
        })
        .error(function(data, status) {
                console.log('Error status:' + status);
        });
		
	$http.get('http://192.168.248.128:8080/data/training')
        .success(function(data) {
            $scope.trainingdataset = data;
			console.log(data)
        })
        .error(function(data, status) {
                console.log('Error status:' + status);
        });
		
	$http.get('http://192.168.248.128:8080/data/validation')
        .success(function(data) {
            $scope.validationdataset = data;
			console.log(data)
        })
        .error(function(data, status) {
                console.log('Error status:' + status);
        });

    // Cuando se añade un nuevo TODO, manda el texto a la API
    /*$scope.createData = function(){
        $http.post('/api/todos', $scope.formData)
            .success(function(data) {
                $scope.formData = {};
                $scope.data = data;
                console.log(data);
            })
            .error(function(data, status) {
                console.log('Error status:' + status);
            });
    };*/
}
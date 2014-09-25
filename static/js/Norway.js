      google.load("visualization", "1", {packages:["geochart"]});
      google.setOnLoadCallback(drawRegionsMap);

      function drawRegionsMap() {

        var data = google.visualization.arrayToDataTable([
          ['Country', 'Display'],
          ['Norway', 1],
        ]);

        var options = {
		legend: 'none',
	        // https://google-developers.appspot.com/chart/interactive/docs/gallery/geochart#Continent_Hierarchy
		region: '154',
		colorAxis: {minValue: 0, colors: ['black', '#dd0000']},
	};

        var chart = new google.visualization.GeoChart(document.getElementById('Norway'));

        chart.draw(data, options);
      }


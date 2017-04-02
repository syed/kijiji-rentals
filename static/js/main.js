$(function () {
    var form = $("#searchform");
    form.submit(function (e) {
        //prevent double submit
        //make spinner visible and hide the search button
        $("#submitbutton input").css("visibility", "hidden");
        $("#submitbutton i").css("visibility", "visible");

        var query_data = form.serializeArray();
        console.log(query_data);

        $.getJSON(form.attr('action'), //url
            query_data,//post data
            function (responseText, responseStatus) {
                console.log(responseText);

                $("#submitbutton input").css("visibility", "visible");
                $("#submitbutton i").css("visibility", "hidden");

                addMarkers(responseText);
            }
        );
        // prevent actual submit from happening
        e.preventDefault();
    });
});

$(function () {
    $("#datepicker").datepicker({minDate: -20, maxDate: "+1M +10D"});
});
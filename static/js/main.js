$(function () {
    var form = $("#searchform");
    form.submit(function (e) {
        //prevent double submit
        $("#submitbutton").attr('disabled', true);

        var query_data = form.serializeArray();
        console.log(query_data);

        $.getJSON(form.attr('action'), //url
            query_data,//post data
            function (responseText, responseStatus) {
                console.log(responseText);

                $("#submitbutton").attr('disabled', false);
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
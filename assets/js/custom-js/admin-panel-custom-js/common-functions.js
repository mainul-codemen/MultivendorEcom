function statusDd(obj) {
    $("#UdStatus").empty().val(obj.Form.Status);
    if (obj.Form.Status == 1) {
        $("#UdStatus").append('<option value="' + 1 + '">' + "Active" + '</option>');
        $("#UdStatus").append('<option value="' + 2 + '">' + "Inactive" + '</option>');
    } else {
        $("#UdStatus").append('<option value="' + 2 + '">' + "Inactive" + '</option>');
        $("#UdStatus").append('<option value="' + 1 + '">' + "Active" + '</option>');
    }
}
// view : Hub Form with Dropdown
function viewHubForm() {
    $.ajax({
        url: '/admin/hub/create',
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data)
            var countries = obj.CountryData
            var cntryDdd = $("#countrydd");
            $("#countrydd").append('<option>--Select Country--</option>');
            $(countries).each(function () {
                var option = $("<option />");
                option.html(this.Name);
                option.val(this.ID);
                cntryDdd.append(option);
            });
            var districts = obj.DistrictData
            var dis = $("#districtdd");
            $("#districtdd").append('<option>--Select District--</option>');
            $(districts).each(function () {
                var option = $("<option />");
                option.html(this.Name);
                option.val(this.ID);
                dis.append(option);
            });

            var stations = obj.StationData
            var stn = $("#stationdd");
            $("#stationdd").append('<option>--Select Station--</option>');
            $(stations).each(function () {
                var option = $("<option />");
                option.html(this.Name);
                option.val(this.ID);
                stn.append(option);
            });
        }
    });
    resetData()
}
// create
$(document).ready(function () {
    $('#saveForm').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        $.ajax({
            url: "/admin/hub/create",
            method: 'post',
            data: $('form.tagForm').serialize(),
            success: function (data) {
                var obj = jQuery.parseJSON(data);
                var Toast = Swal.mixin({
                    toast: true,
                    position: 'top-end',
                    showConfirmButton: false,
                    timer: 3000
                });
                if (obj.Status) {
                    $('#saveForm').trigger("reset");
                    Toast.fire({
                        icon: 'success',
                        title: obj.Message
                    })
                    setTimeout(function () { window.location.reload(true); }, 1000);
                    $('#modal-add').modal('hide');
                } else {
                    $("#Name").empty().append(obj.FormErrors.Name);
                    $("#CountryID").empty().append(obj.FormErrors.CountryID);
                    $("#DistrictID").empty().append(obj.FormErrors.DistrictID);
                    $("#StationID").empty().append(obj.FormErrors.StationID);
                    $("#HubPhone1").empty().append(obj.FormErrors.HubPhone1);
                    $("#HubPhone2").empty().append(obj.FormErrors.HubPhone2);
                    $("#HubEmail").empty().append(obj.FormErrors.HubEmail);
                    $("#HubAddress").empty().append(obj.FormErrors.HubAddress);
                    $("#Position").empty().append(obj.FormErrors.Position);
                    $("#Status").empty().append(obj.FormErrors.Status);
                    Toast.fire({
                        icon: 'error',
                        title: "Please Insert All Data Carefully."
                    })
                }
            },
        });
        // reset all form data after close modal
        $('#modal-add').on('hidden.bs.modal', function () {
            $(this).find('form').trigger('reset');
            $("#Name").empty();
            $("#CountryID").empty();
            $("#DistrictID").empty();
            $("#StationID").empty();
            $("#HubPhone1").empty();
            $("#HubPhone2").empty();
            $("#HubEmail").empty();
            $("#HubAddress").empty();
            $("#Position").empty();
            $("#Status").empty();
        });
    });
});

// view : Hub
function viewHub(id) {
    $.ajax({
        url: "/admin/hub/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#VName").empty().append(obj.Form.Name);
            $("#VCountryName").empty().append(obj.Form.CountryName);
            $("#VDistrictName").empty().append(obj.Form.DistrictName);
            $("#VStationName").empty().append(obj.Form.StationName);
            $("#VHubPhone1").empty().append(obj.Form.HubPhone1);
            $("#VHubPhone2").empty().append(obj.Form.HubPhone2);
            $("#VHubEmail").empty().append(obj.Form.HubEmail);
            $("#VHubAddress").empty().append(obj.Form.HubAddress);
            $("#VPosition").empty().append(obj.Form.Position);
            if (obj.Form.Status == 1) {
                $("#VStatus").empty().append("Active");
            }else{
                $("#VStatus").empty().append("InActive");
            }
        }
    })
}

// Update : Hub View
function viewHubUpdateData(id) {
    $.ajax({
        url: "/admin/hub/update/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#UdID").empty().val(obj.Form.ID);
            $("#UdName").empty().val(obj.Form.Name);
            $("#UdHubPhone1").empty().val(obj.Form.HubPhone1);
            $("#UdHubPhone2").empty().val(obj.Form.HubPhone2);
            $("#UdHubEmail").empty().val(obj.Form.HubEmail);
            $("#UdHubAddress").empty().val(obj.Form.HubAddress);
            $("#UdPosition").empty().val(obj.Form.Position);
            $("#UdStatus").empty().val(obj.Form.Status);
            hubDropdownUpdate(obj);
        }
    });
    resetDataUpdate();
}
// Update : Hub Submit
$(document).ready(function () {
    $('#updateForm').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#UdID').val();
        $.ajax({
            url: "/admin/hub/update/" + id,
            method: 'post',
            data: $('form.tagUpForm').serialize(),
            success: function (data) {
                var obj = jQuery.parseJSON(data);
                var Toast = Swal.mixin({
                    toast: true,
                    position: 'top-end',
                    showConfirmButton: false,
                    timer: 3000
                });
                if (obj.Status) {
                    $('updateForm').trigger("reset");
                    Toast.fire({
                        icon: 'success',
                        title: obj.Message
                    })
                    setTimeout(function () { window.location.reload(true); }, 1000);
                    $('#modal-update').modal('hide');
                } else {
                    $("#NameErr").empty().append(obj.FormErrors.Name);
                    $("#CountryIDErr").empty().append(obj.FormErrors.CountryID);
                    $("#DistrictIDErr").empty().append(obj.FormErrors.DistrictID);
                    $("#StationIDErr").empty().append(obj.FormErrors.StationID);
                    $("#HubPhone1Err").empty().append(obj.FormErrors.HubPhone1);
                    $("#HubPhone2Err").empty().append(obj.FormErrors.HubPhone2);
                    $("#HubEmailErr").empty().append(obj.FormErrors.HubEmail);
                    $("#HubAddressErr").empty().append(obj.FormErrors.HubAddress);
                    $("#PositionErr").empty().append(obj.FormErrors.Position);
                    $("#StatusErr").empty().append(obj.FormErrors.Status);
                    Toast.fire({
                        icon: 'error',
                        title: "Please Insert All Data Carefully."
                    })
                }
            }
        });
        // reset all form data after close modal
        $('#modal-update').on('hidden.bs.modal', function () {
            $(this).find('form').trigger('reset');
            $("#NameErr").empty();
            $("#CountryIDErr").empty();
            $("#DistrictIDErr").empty();
            $("#StationIDErr").empty();
            $("#HubPhone1Err").empty();
            $("#HubPhone2Err").empty();
            $("#HubEmailErr").empty();
            $("#HubAddressErr").empty();
            $("#PositionErr").empty();
            $("#StatusErr").empty();
        });
    });
});
// Delete : Country View
function deleteHubData(id) {
    $.ajax({
        url: "/admin/hub/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#dID").empty().val(obj.Form.ID);
            $("#dName").empty().append(obj.Form.Name);
            $("#dCountryName").empty().append(obj.Form.CountryName);
            $("#dDistrictName").empty().append(obj.Form.DistrictName);
            $("#dStationName").empty().append(obj.Form.StationName);
            $("#dPosition").empty().append(obj.Form.Position);
            $("#dStatus").empty().append(obj.Form.Status);
            if (obj.Form.Status == 1) {
                $("#dStatus").empty().append("Active");
            }else{
                $("#dStatus").empty().append("InActive");
            }
        }
    })
}
// delete hub
$(document).ready(function () {
    $('#deleteHub').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#dID').val();
        $.ajax({
            url: "/admin/hub/delete/" + id,
            method: 'get',
            success: function (data) {
                var obj = jQuery.parseJSON(data);
                var Toast = Swal.mixin({
                    toast: true,
                    position: 'top-end',
                    showConfirmButton: false,
                    timer: 3000
                });
                if (obj.Status) {
                    $('deleteCountry').trigger("reset");
                    Toast.fire({
                        icon: 'success',
                        title: obj.Message
                    })
                    setTimeout(function () { window.location.reload(true); }, 1000);
                    $('#modal-delete').modal('hide');
                }
            }
        });
    });
});

// dropdown Update
function hubDropdownUpdate(obj) {
    var countries = obj.CountryData
    var cntryDdd = $("#countrydd-update");
    $("#countrydd-update").append('<option>--Select Country--</option>');
    $(countries).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.ID);
        cntryDdd.append(option);
    });
    var districts = obj.DistrictData
    var dis = $("#districtdd-update");
    $("#districtdd-update").append('<option>--Select District--</option>');
    $(districts).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.ID);
        dis.append(option);
    });

    var stations = obj.StationData
    var stn = $("#stationdd-update");
    $("#stationdd-update").append('<option>--Select Station--</option>');
    $(stations).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.ID);
        stn.append(option);
    });
}
// reset dropdown data
function resetData() {
    // reset all form data after close modal
    $('#modal-add').on('hidden.bs.modal', function () {
        $(this).find('form').trigger('reset');
        $("#countrydd").empty();
        $("#districtdd").empty();
        $("#stationdd").empty();
    });
}
// reset dropdown data update
function resetDataUpdate() {
    // reset all form data after close modal
    $('#modal-update').on('hidden.bs.modal', function () {
        $(this).find('form').trigger('reset');
        $("#countrydd-update").empty();
        $("#districtdd-update").empty();
        $("#stationdd-update").empty();
    });
}

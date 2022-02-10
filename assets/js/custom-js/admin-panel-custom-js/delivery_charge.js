// view : DeliveryCharge Form with Dropdown
function viewDeliveryChargeForm() {
    $.ajax({
        url: '/admin/delivery-charge/create',
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            dcDropdown(obj)
        }
    });
}


// dropdown
function dcDropdown(obj) {
    var districts = obj.DistrictData
    var dis = $("#districtdd");
    $(districts).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.DistrictID);
        dis.append(option);
    });
    var cntrys = obj.CountryData
    var cntry = $("#countrydd");
    $(cntrys).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.CountryID);
        cntry.append(option);
    });
    var stations = obj.StationData
    var stn = $("#stationdd");
    $(stations).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.StationID);
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
// create
$(document).ready(function () {
    $('#saveForm').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        $.ajax({
            url: "/admin/delivery-company/create",
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
                    $("#CountryID").empty().append(obj.FormErrors.CountryID);
                    $("#DistrictID").empty().append(obj.FormErrors.DistrictID);
                    $("#StationID").empty().append(obj.FormErrors.StationID);
                    $("#MinWeight").empty().append(obj.FormErrors.MinWeight);
                    $("#MaxWeight").empty().append(obj.FormErrors.MaxWeight);
                    $("#DeliveryCharge").empty().append(obj.FormErrors.DeliveryCharge);
                    $("#DCStatus").empty().append(obj.FormErrors.DCStatus);
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
            $("#CountryID").empty();
            $("#DistrictID").empty();
            $("#StationID").empty();
            $("#MinWeight").empty();
            $("#MaxWeight").empty();
            $("#DeliveryCharge").empty();
            $("#DCStatus").empty();
        });
    });
});

// view : DeliveryCharge
function viewDeliveryCharge(id) {
    $.ajax({
        url: "/admin/delivery-charge/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#VCountryName").empty().append(obj.Form.CountryName);
            $("#VDistrictName").empty().append(obj.Form.DistrictName);
            $("#VStationName").empty().append(obj.Form.StationName);
            $("#VMinWeight").empty().append(obj.Form.MinWeight);
            $("#VMaxWeight").empty().append(obj.Form.MaxWeight);
            $("#VDeliveryCharge").empty().append(obj.Form.DeliveryCharge);
            $("#VDCStatus").empty().append(obj.Form.DCStatus);
        }
    })
}

// Update : DeliveryCharge View
function viewDeliveryChargeUpdateData(id) {
    $.ajax({
        url: "/admin/delivery-company/update/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#UdID").empty().val(obj.Form.ID);
            $("#UdPhone").empty().val(obj.Form.Phone);
            $("#UdEmail").empty().val(obj.Form.Email);
            $("#UdCompanyAddress").empty().val(obj.Form.CompanyAddress);
            $("#UdPosition").empty().val(obj.Form.Position);
            $("#UdCompanyStatus").empty().val(obj.Form.CompanyStatus);
            dcDropdownUpdate(obj)
        }
    });
    resetDataUpdate()
}

// dropdown Update
function dcDropdownUpdate(obj) {
    var districts = obj.DistrictData
    var dis = $("#districtdd-update");
    $(districts).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.DistrictID);
        dis.append(option);
    });
    var cntrys = obj.CountryData
    var cntry = $("#countrydd-update");
    $(cntrys).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.CountryID);
        cntry.append(option);
    });
    var stations = obj.StationData
    var stn = $("#stationdd-update");
    $(stations).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.StationID);
        stn.append(option);
    });
}
// Update : DeliveryCharge Submit
$(document).ready(function () {
    $('#updateForm').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#UdID').val();
        $.ajax({
            url: "/admin/delivery-company/update/" + id,
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
                    $("#CountryIDErr").empty().append(obj.FormErrors.CountryID);
                    $("#DistrictIDErr").empty().append(obj.FormErrors.DistrictID);
                    $("#StationIDErr").empty().append(obj.FormErrors.StationID);
                    $("#MainWeightErr").empty().append(obj.FormErrors.MinWeight);
                    $("#MaxWeightErr").empty().append(obj.FormErrors.MaxWeight);
                    $("#DeliveryChargeErr").empty().append(obj.FormErrors.DeliveryCharge);
                    $("#DCStatusErr").empty().append(obj.FormErrors.DCStatus);
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
            $("#CountryIDErr").empty();
            $("#DistrictIDErr").empty();
            $("#StationIDErr").empty();
            $("#MinWeightErr").empty();
            $("#MaxWeightErr").empty();
            $("#DeliveryChargeErr").empty();
            $("#DCStatusErr").empty();
        });
    });
});
// Delete : Country View
function deleteDeliveryChargeData(id) {
    $.ajax({
        url: "/admin/delivery-company/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#dID").empty().val(obj.Form.ID);
            $("#dCompanyName").empty().append(obj.Form.Name);
            $("#dCountryName").empty().append(obj.Form.CountryName);
            $("#dDistrictName").empty().append(obj.Form.DistrictName);
            $("#dStationName").empty().append(obj.Form.StationName);
            $("#dStatus").empty().append(obj.Form.DCStatus);
        }
    })
}
// delete country
$(document).ready(function () {
    $('#deleteDeliveryCharge').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#dID').val();
        $.ajax({
            url: "/admin/delivery-company/delete/" + id,
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

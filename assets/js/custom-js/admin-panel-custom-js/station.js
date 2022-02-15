
// view : Station Form with Dropdown district
function viewStationForm() {
    $.ajax({
        url: '/admin/station/create',
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data)
            var countries = obj.DistrictData
            var cntryDdd = $("#districtdd");
            $("#districtdd").append('<option value="">--Select District--</option>');
            $(countries).each(function () {
                var option = $("<option />");
                option.html(this.Name);
                option.val(this.ID);
                cntryDdd.append(option);
            });
        }
    });
    // reset all form data after close modal
    $('#modal-add').on('hidden.bs.modal', function () {
        $(this).find('form').trigger('reset');
        $("#districtdd").empty();
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
            url: "/admin/station/create",
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
                    $("#DistrictID").empty().append(obj.FormErrors.DistrictID);
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
            $("#DistrictID").empty();
            $("#Position").empty();
            $("#Status").empty();
        });
    });
});

// view : Station
function viewStation(id) {
    $.ajax({
        url: "/admin/station/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#VName").empty().append(obj.Form.Name);
            $("#VDistrictName").empty().append(obj.Form.DistrictName);
            $("#VPosition").empty().append(obj.Form.Position);
            $("#VStatus").empty().append(obj.Form.Status);
            if (obj.Form.Status == 1) {
                $("#VStatus").empty().append("Active");
            }else{
                $("#VStatus").empty().append("InActive");
            }
        }
    })
}

// Update : Station View
function viewStationUpdateData(id) {
    $.ajax({
        url: "/admin/station/update/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#UdID").empty().val(obj.Form.ID);
            $("#UdName").empty().val(obj.Form.Name);
            $("#UdPosition").empty().val(obj.Form.Position);
            $("#UdStatus").empty().val(obj.Form.Status);
            var districts = obj.DistrictData
            var districtdd = $("#districtdd-update");
            $("#districtdd-update").append('<option value="">--Select District--</option>');
            $(districts).each(function () {
                var option = $("<option />");
                option.html(this.Name);
                option.val(this.ID);
                districtdd.append(option);
            });
        }
    });
    // reset all form data after close modal
    $('#modal-update').on('hidden.bs.modal', function () {
        $(this).find('form').trigger('reset');
        $("#districtdd-update").empty();
    });
}
// Update : Station Submit
$(document).ready(function () {
    $('#updateForm').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#UdID').val();
        $.ajax({
            url: "/admin/station/update/" + id,
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
                    $("#DistrictIDErr").empty().append(obj.FormErrors.DistrictID);
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
            $("#DistrictErr").empty();
            $("#PositionErr").empty();
            $("#StatusErr").empty();
        });
    });
});
// Delete : District View
function deleteStationData(id) {
    $.ajax({
        url: "/admin/station/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#dID").empty().val(obj.Form.ID);
            $("#dName").empty().append(obj.Form.Name);
            $("#dPosition").empty().append(obj.Form.Position);
            if (obj.Form.Status == 1) {
                $("#dStatus").empty().append("Active");
            }else{
                $("#dStatus").empty().append("InActive");
            }
        }
    })
}
// delete country
$(document).ready(function () {
    $('#deleteStation').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#dID').val();
        $.ajax({
            url: "/admin/station/delete/" + id,
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
                    $('deleteDistrict').trigger("reset");
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

// create
$(document).ready(function () {
    $('#saveForm').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        $.ajax({
            url: "/admin/district/create",
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
            $("#Position").empty();
            $("#Status").empty();
        });
    });
});
// view : District Form with Dropdown
function viewDistrictForm() {
    $.ajax({
        url: '/admin/district/create',
        method: 'get',
        success: function (data) {
            console.log(data)
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
        }
    });
    // reset all form data after close modal
    $('#modal-add').on('hidden.bs.modal', function () {
        $(this).find('form').trigger('reset');
        $("#countrydd").empty();
    });
}
// view : District
function viewDistrict(id) {
    $.ajax({
        url: "/admin/district/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#VName").empty().append(obj.Form.Name);
            $("#VCountryName").empty().append(obj.Form.CountryName);
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

// Update : District View with Country Dropdown
function viewDistrictUpdateData(id) {
    $.ajax({
        url: "/admin/district/update/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#UdID").empty().val(obj.Form.ID);
            $("#UdName").empty().val(obj.Form.Name);
            $("#UdPosition").empty().val(obj.Form.Position);
            statusDd(obj);
            var countries = obj.CountryData
            var cntdd = $("#countrydd-update");
            $("#countrydd-update").append('<option value="'+obj.Form.CountryID+'">'+obj.Form.CountryName+'</option>');
            $(countries).each(function () {
                if (this.Name != obj.Form.CountryName){
                    var option = $("<option />");
                    option.html(this.Name);
                    option.val(this.ID);
                    cntdd.append(option);
                }
            });
        }
    });
    // reset all form data after close modal
    $('#modal-update').on('hidden.bs.modal', function () {
        $(this).find('form').trigger('reset');
        $("#countrydd-update").empty();
    });
}

// Update : District Submit
$(document).ready(function () {
    $('#updateForm').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#UdID').val();
        $.ajax({
            url: "/admin/district/update/" + id,
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
            $("#CountryErr").empty();
            $("#PositionErr").empty();
            $("#StatusErr").empty();
        });
    });
});

// Delete : Country View
function deleteDistrictData(id) {
    $.ajax({
        url: "/admin/district/view/" + id,
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
    $('#deleteDistrict').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#dID').val();
        $.ajax({
            url: "/admin/district/delete/" + id,
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

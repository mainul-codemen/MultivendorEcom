// create : UserRole
$(document).ready(function () {
    $('#saveForm').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        $.ajax({
            url: "/admin/user-role/create",
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
                    $("#Description").empty().append(obj.FormErrors.Description);
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
            $("#Description").empty();
            $("#Position").empty();
            $("#Status").empty();
        });
    });
});
// view : UserRole
function viewUserRole(id) {
    $.ajax({
        url: "/admin/user-role/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#VName").empty().append(obj.Form.Name);
            $("#VDescription").empty().append(obj.Form.Description);
            $("#VPosition").empty().append(obj.Form.Position);
            if (obj.Form.Status == 1) {
                $("#VStatus").empty().append("Active");
            } else {
                $("#VStatus").empty().append("InActive");
            }
        }
    })
}
// Update : UserRole View
function viewUserRoleUpdateData(id) {
    $.ajax({
        url: "/admin/user-role/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#UdID").empty().val(obj.Form.ID);
            $("#UdName").empty().val(obj.Form.Name);
            $("#UdDescription").empty().val(obj.Form.Description);
            $("#UdPosition").empty().val(obj.Form.Position);
            $("#UdStatus").empty().val(obj.Form.Status);
            statusDd(obj);
        }
    })
}
// Update : UserRole Submit
$(document).ready(function () {
    $('#updateForm').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#UdID').val();
        $.ajax({
            url: "/admin/user-role/update/" + id,
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
                    $("#DescriptionErr").empty().append(obj.FormErrors.Description);
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
            $("#DescriptionErr").empty();
            $("#PositionErr").empty();
            $("#StatusErr").empty();
        });
    });
});

// Delete : UserRole View
function deleteUserRoleData(id) {
    $.ajax({
        url: "/admin/user-role/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#dID").empty().val(obj.Form.ID);
            $("#dName").empty().append(obj.Form.Name);
            $("#dDescription").empty().append(obj.Form.Description);
            $("#dPosition").empty().append(obj.Form.Position);
            if (obj.Form.Status == 1) {
                $("#dStatus").empty().append("Active");
            } else {
                $("#dStatus").empty().append("InActive");
            }
        }
    })
}
// delete designaiton
$(document).ready(function () {
    $('#deleteUserRole').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#dID').val();
        $.ajax({
            url: "/admin/user-role/delete/" + id,
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
                    $('deleteUserRole').trigger("reset");
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

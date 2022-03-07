// create : TransactionTypes
$(document).ready(function () {
    $('#saveForm').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        $.ajax({
            url: "/admin/transaction-types/create",
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
                    $("#TransactionTypesName").empty().append(obj.FormErrors.TransactionTypesName);
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
            $("#TransactionTypesName").empty();
            $("#Status").empty();
        });
    });
});

// view : TransactionTypes
function viewTransactionTypes(id) {
    $.ajax({
        url: "/admin/transaction-types/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#VTransactionTypesName").empty().append(obj.Form.TransactionTypesName);
            if (obj.Form.Status == 1) {
                $("#VStatus").empty().append("Active");
            } else {
                $("#VStatus").empty().append("InActive");
            }
        }
    });
}

// Update : TransactionTypes View
function viewTransactionTypesUpdateData(id) {
    $.ajax({
        url: "/admin/transaction-types/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#UdID").empty().val(obj.Form.ID);
            $("#UdTransactionTypesName").empty().val(obj.Form.TransactionTypesName);
            statusDd(obj);
        }
    });
}

// Update : TransactionTypes Submit
$(document).ready(function () {
    $('#updateForm').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#UdID').val();
        $.ajax({
            url: "/admin/transaction-types/update/" + id,
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
                    $("#TransactionTypesNameErr").empty().append(obj.FormErrors.TransactionTypesName);
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
            $("#TransactionTypesNameErr").empty();
            $("#StatusErr").empty();
        });
    });
});

// Delete : TransactionTypes View
function deleteTransactionTypesData(id) {
    $.ajax({
        url: "/admin/transaction-types/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#dID").empty().val(obj.Form.ID);
            $("#dTransactionTypesName").empty().append(obj.Form.TransactionTypesName);
            if (obj.Form.Status == 1) {
                $("#dStatus").empty().append("Active");
            } else {
                $("#dStatus").empty().append("InActive");
            }
        }
    });
}

// delete designaiton
$(document).ready(function () {
    $('#deleteTransactionTypes').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#dID').val();
        $.ajax({
            url: "/admin/transaction-types/delete/" + id,
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
                    $('deleteTransactionTypes').trigger("reset");
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
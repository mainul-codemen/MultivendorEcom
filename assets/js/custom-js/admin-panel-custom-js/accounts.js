// create : Accounts
$(document).ready(function () {
    $('#saveForm').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        $.ajax({
            url: "/admin/accounts/create",
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
                    $("#AccountName").empty().append(obj.FormErrors.AccountName);
                    $("#AccountVisualization").empty().append(obj.FormErrors.AccountVisualization);
                    $("#AccountNumber").empty().append(obj.FormErrors.AccountNumber);
                    $("#Amount").empty().append(obj.FormErrors.Amount);
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
            $("#AccountName").empty();
            $("#AccountVisualization").empty();
            $("#AccountNumber").empty();
            $("#Amount").empty();
            $("#Status").empty();
        });
    });
});

// view : Accounts
function viewAccounts(id) {
    $.ajax({
        url: "/admin/accounts/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#VAccountName").empty().append(obj.Form.AccountName);
            $("#VAccountNumber").empty().append(obj.Form.AccountNumber);
            $("#VAmount").empty().append(obj.Form.Amount);
            if (obj.Form.Status == 1) {
                $("#VStatus").empty().append("Active");
            }else{
                $("#VStatus").empty().append("InActive");
            }
            if (obj.Form.AccountVisualization == 1) {
                $("#VAccountVisualization").empty().append("Front Visualization");
            }else{
                $("#VAccountVisualization").empty().append("Visualization Restriction");
            }
        }
    })
}

// Update : Accounts View
function viewAccountsUpdateData(id) {
    $.ajax({
        url: "/admin/accounts/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#UdID").empty().val(obj.Form.ID);
            $("#UdAccountName").empty().val(obj.Form.AccountName);
            $("#UdAccountVisualization").empty().val(obj.Form.AccountVisualization);
            $("#UdAccountNumber").empty().val(obj.Form.AccountNumber);
            $("#UdAmount").empty().val(obj.Form.Amount);
            statusDb(obj)
            accountVisual(obj)
        }
    });
}

function accountVisual(obj) {
    $("#UdAccountVisualization").empty().val(obj.Form.AccountVisualization);
    if (obj.Form.AccountVisualization == 1) {
        $("#UdAccountVisualization").append('<option value="' + 1 + '">' + "Front Visualization" + '</option>');
        $("#UdAccountVisualization").append('<option value="' + 2 + '">' + "Visualization Restriction" + '</option>');
    } else {
        $("#UdAccountVisualization").append('<option value="' + 2 + '">' + "Visualization Restriction" + '</option>');
        $("#UdAccountVisualization").append('<option value="' + 1 + '">' + "Front Visualization" + '</option>');
    }
}

function statusDb(obj) {
    $("#UdStatus").empty().val(obj.Form.Status);
    if (obj.Form.Status == 1) {
        $("#UdStatus").append('<option value="' + 1 + '">' + "Active" + '</option>');
        $("#UdStatus").append('<option value="' + 2 + '">' + "Inactive" + '</option>');
    } else {
        $("#UdStatus").append('<option value="' + 2 + '">' + "Inactive" + '</option>');
        $("#UdStatus").append('<option value="' + 1 + '">' + "Active" + '</option>');
    }
}

// Update : Accounts Submit
$(document).ready(function () {
    $('#updateForm').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#UdID').val();
        $.ajax({
            url: "/admin/accounts/update/" + id,
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
                    $("#AccountNameErr").empty().append(obj.FormErrors.AccountName);
                    $("#AccountNumberErr").empty().append(obj.FormErrors.AccountNumber);
                    $("#AccountVisualizationErr").empty().append(obj.FormErrors.AccountVisualization);
                    $("#AmountErr").empty().append(obj.FormErrors.Amount);
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
            $("#AccountNameErr").empty();
            $("#AccountNumberErr").empty();
            $("#AccountVisualizationErr").empty();
            $("#AmountErr").empty();
            $("#StatusErr").empty();
        });
    });
});

// Delete : Accounts View
function deleteAccountsData(id) {
    $.ajax({
        url: "/admin/accounts/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#dID").empty().val(obj.Form.ID);
            $("#dAccountName").empty().append(obj.Form.AccountName);
            $("#dAccountNumber").empty().append(obj.Form.AccountNumber);
            $("#dAmount").empty().append(obj.Form.Amount);
            if (obj.Form.AccountVisualization == 1) {
                $("#dAccountVisualization").empty().append("Front Visualization");
            }else{
                $("#dAccountVisualization").empty().append("Visualization Restriction");
            }
            if (obj.Form.Status == 1) {
                $("#dStatus").empty().append("Active");
            }else{
                $("#dStatus").empty().append("InActive");
            }
        }
    })
}

// delete accounts
$(document).ready(function () {
    $('#deleteAccounts').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#dID').val();
        $.ajax({
            url: "/admin/accounts/delete/" + id,
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
                    $('deleteAccounts').trigger("reset");
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

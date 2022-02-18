// create : Grade
$(document).ready(function () {
    $('#saveForm').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        $.ajax({
            url: "/admin/grade/create",
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
                    $("#BasicSalary").empty().append(obj.FormErrors.BasicSalary);
                    $("#LunchAllowance").empty().append(obj.FormErrors.LunchAllowance);
                    $("#Transportation").empty().append(obj.FormErrors.Transportation);
                    $("#RentAllowance").empty().append(obj.FormErrors.RentAllowance);
                    $("#AbsentPenalty").empty().append(obj.FormErrors.AbsentPenalty);
                    $("#TotalSalary").empty().append(obj.FormErrors.TotalSalary);
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
            $("#BasicSalary").empty();
            $("#LunchAllowance").empty();
            $("#Transportation").empty();
            $("#RentAllowance").empty();
            $("#AbsentPenalty").empty();
            $("#TotalSalary").empty();
            $("#Position").empty();
            $("#Status").empty();
        });
    });
});
// view : Grade
function viewGrade(id) {
    $.ajax({
        url: "/admin/grade/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#VName").empty().append(obj.Form.Name);
            $("#VBasicSalary").empty().append(obj.Form.BasicSalary);
            $("#VLunchAllowance").empty().append(obj.Form.LunchAllowance);
            $("#VTransportation").empty().append(obj.Form.Transportation);
            $("#VRentAllowance").empty().append(obj.Form.RentAllowance);
            $("#VAbsentPenalty").empty().append(obj.Form.AbsentPenalty);
            $("#VTotalSalary").empty().append(obj.Form.TotalSalary);
            $("#VStatus").empty().append(obj.Form.Status);
            $("#VPosition").empty().append(obj.Form.Position);
            if (obj.Form.Status == 1) {
                $("#VStatus").empty().append("Active");
            } else {
                $("#VStatus").empty().append("InActive");
            }
        }
    })
}
// Update : Grade View
function viewGradeUpdateData(id) {
    $.ajax({
        url: "/admin/grade/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#UdID").empty().val(obj.Form.ID);
            $("#UName").empty().append(obj.Form.Name);
            $("#UBasicSalary").empty().append(obj.Form.BasicSalary);
            $("#ULunchAllowance").empty().append(obj.Form.LunchAllowance);
            $("#UTransportation").empty().append(obj.Form.Transportation);
            $("#URentAllowance").empty().append(obj.Form.RentAllowance);
            $("#UAbsentPenalty").empty().append(obj.Form.AbsentPenalty);
            $("#UTotalSalary").empty().append(obj.Form.TotalSalary);
            $("#UStatus").empty().append(obj.Form.Status);
            $("#UPosition").empty().append(obj.Form.Position);
            statusDd(obj);
        }
    })
}
// Update : Grade Submit
$(document).ready(function () {
    $('#updateForm').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#UdID').val();
        $.ajax({
            url: "/admin/grade/update/" + id,
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
                    $("#BasicSalaryErr").empty().append(obj.FormErrors.BasicSalary);
                    $("#TransportationErr").empty().append(obj.FormErrors.Transportation);
                    $("#LunchAllowanceErr").empty().append(obj.FormErrors.LunchAllowance);
                    $("#RentAllowanceErr").empty().append(obj.FormErrors.RentAllowance);
                    $("#AbsentPenaltyErr").empty().append(obj.FormErrors.AbsentPenalty);
                    $("#TotalSalaryErr").empty().append(obj.FormErrors.TotalSalary);
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
            $("#BasicSalaryErr").empty();
            $("#LunchAllowanceErr");
            $("#RentAllowanceErr");
            $("#TransportationErr");
            $("#AbsentPenaltyErr");
            $("#TotalSalaryErr");
            $("#PositionErr");
            $("#StatusErr");
        });
    });
});

// Delete : Grade View
function deleteGradeData(id) {
    $.ajax({
        url: "/admin/grade/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#dID").empty().val(obj.Form.ID);
            $("#dName").empty().append(obj.Form.Name);
            $("#dBasicSalary").empty().append(obj.Form.BasicSalary);
            $("#dLunchAllowance").empty().append(obj.Form.LunchAllowance);
            $("#dTransportation").empty().append(obj.Form.Transportation);
            $("#dRentAllowance").empty().append(obj.Form.RentAllowance);
            $("#dAbsentPenalty").empty().append(obj.Form.AbsentPenalty);
            $("#dTotalSalary").empty().append(obj.Form.TotalSalary);
            
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
    $('#deleteGrade').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#dID').val();
        $.ajax({
            url: "/admin/grade/delete/" + id,
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
                    $('deleteGrade').trigger("reset");
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

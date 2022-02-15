// view : Branch Form with Dropdown
function viewBranchForm() {
    $.ajax({
        url: '/admin/branch/create',
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            branchDropdown(obj)
        }
    });
    resetData()
}

// dropdown
function branchDropdown(obj) {
    var cntrys = obj.CountryData
    var cntrydd = $("#countrydd");
    $("#countrydd").append('<option>--Select Country--</option>');
    $(cntrys).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.ID);
        cntrydd.append(option);
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
    var stndd = $("#stationdd");
    $("#stationdd").append('<option>--Select Station--</option>');
    $(stations).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.ID);
        stndd.append(option);
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
            url: "/admin/branch/create",
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
                    $("#BranchPhone1").empty().append(obj.FormErrors.BranchPhone1);
                    $("#BranchPhone2").empty().append(obj.FormErrors.BranchPhone2);
                    $("#BranchEmail").empty().append(obj.FormErrors.BranchEmail);
                    $("#BranchAddress").empty().append(obj.FormErrors.BranchAddress);
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
            $("#BranchPhone1").empty();
            $("#BranchPhone2").empty();
            $("#BranchEmail").empty();
            $("#BranchAddress").empty();
            $("#Position").empty();
            $("#Status").empty();
        });
    });
});

// view : Branch
function viewBranch(id) {
    $.ajax({
        url: "/admin/branch/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#VName").empty().append(obj.Form.Name);
            $("#VCountryName").empty().append(obj.Form.CountryName);
            $("#VDistrictName").empty().append(obj.Form.DistrictName);
            $("#VStationName").empty().append(obj.Form.StationName);
            $("#VBranchPhone1").empty().append(obj.Form.BranchPhone1);
            $("#VBranchPhone2").empty().append(obj.Form.BranchPhone2);
            $("#VBranchEmail").empty().append(obj.Form.BranchEmail);
            $("#VBranchAddress").empty().append(obj.Form.BranchAddress);
            $("#VPosition").empty().append(obj.Form.Position);
            if (obj.Form.BranchStatus == 1) {
                $("#VStatus").empty().append("Active");
            } else {
                $("#VStatus").empty().append("InActive");
            }
        }
    })
}

// dropdown Update
function branchDropdownUpdate(obj) {
    var cntrys = obj.CountryData
    var cntrydd = $("#countrydd-update");
    $("#countrydd-update").append('<option>--Select Country--</option>');
    $(cntrys).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.ID);
        cntrydd.append(option);
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
    var stndd = $("#stationdd-update");
    $("#stationdd-update").append('<option>--Select Station--</option>');
    $(stations).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.ID);
        stndd.append(option);
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
// Update : Branch View
function viewBranchUpdateData(id) {
    $.ajax({
        url: "/admin/branch/update/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#UdID").empty().val(obj.Form.ID);
            $("#UdName").empty().val(obj.Form.Name);
            $("#UdBranchPhone1").empty().val(obj.Form.BranchPhone1);
            $("#UdBranchPhone2").empty().val(obj.Form.BranchPhone2);
            $("#UdBranchEmail").empty().val(obj.Form.BranchEmail);
            $("#UdBranchAddress").empty().val(obj.Form.BranchAddress);
            $("#UdPosition").empty().val(obj.Form.Position);
            statusDb(obj);
            branchDropdownUpdate(obj);
        }
    });
    resetDataUpdate()
}

function statusDb(obj) {
    $("#UdStatus").empty().val(obj.Form.BranchStatus);
    if (obj.Form.BranchStatus == 1) {
        $("#UdStatus").append('<option value="' + 1 + '">' + "Active" + '</option>');
        $("#UdStatus").append('<option value="' + 2 + '">' + "Inactive" + '</option>');
    } else {
        $("#UdStatus").append('<option value="' + 2 + '">' + "Inactive" + '</option>');
        $("#UdStatus").append('<option value="' + 1 + '">' + "Active" + '</option>');
    }
}

// Update : Branch Submit
$(document).ready(function () {
    $('#updateForm').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#UdID').val();
        $.ajax({
            url: "/admin/branch/update/" + id,
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
                    $("#BranchPhone1Err").empty().append(obj.FormErrors.BranchPhone1);
                    $("#BranchPhone2Err").empty().append(obj.FormErrors.BranchPhone2);
                    $("#BranchEmailErr").empty().append(obj.FormErrors.BranchEmail);
                    $("#BranchAddressErr").empty().append(obj.FormErrors.BranchAddress);
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
            $("#BranchPhone1Err").empty();
            $("#BranchPhone2Err").empty();
            $("#BranchEmailErr").empty();
            $("#BranchAddressErr").empty();
            $("#PositionErr").empty();
            $("#StatusErr").empty();
        });
    });
});
// Delete : Country View
function deleteBranchData(id) {
    $.ajax({
        url: "/admin/branch/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#dID").empty().val(obj.Form.ID);
            $("#dName").empty().append(obj.Form.Name);
            $("#dCountryName").empty().append(obj.Form.CountryName);
            $("#dDistrictName").empty().append(obj.Form.DistrictName);
            $("#dStationName").empty().append(obj.Form.StationName);
            $("#dPosition").empty().append(obj.Form.Position);
            if (obj.Form.BranchStatus == 1) {
                $("#dStatus").empty().append("Active");
            } else {
                $("#dStatus").empty().append("InActive");
            }
        }
    })
}
// delete country
$(document).ready(function () {
    $('#deleteBranch').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#dID').val();
        $.ajax({
            url: "/admin/branch/delete/" + id,
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

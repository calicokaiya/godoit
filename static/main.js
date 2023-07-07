$(document).ready(function () {
    $('#createTaskForm').submit(function (event) {
        event.preventDefault(); // Prevent the default form submission

        // Serialize the form data
        var formData = $(this).serialize();
        console.log(formData)

        // Send an AJAX request
        $.ajax({
            url: '/tasks/create',
            type: 'POST',
            data: formData,
            success: function (response) {
                // Handle the successful response
                console.log(response);
            },
            error: function (xhr, status, error) {
                // Handle the error
                console.error(error);
            }
        });
    });
    $('#deleteTaskForm').submit(function (event) {
        event.preventDefault(); // Prevent the default form submission

        // Serialize the form data
        var formData = $(this).serialize();

        // Send an AJAX request
        $.ajax({
            url: '/tasks/create',
            type: 'POST',
            data: formData,
            success: function (response) {
                // Handle the successful response
                console.log(response);
            },
            error: function (xhr, status, error) {
                // Handle the error
                console.error(error);
            }
        });
    });
    $('.deleteButton').click(function () {
        var taskId = $(this).data('task-id');
        console.log(taskId);
        var formData = $(this).serialize();

        // Send an AJAX request
        $.ajax({
            url: '/tasks/delete',
            type: 'POST',
            data: "id=" + taskId,
            success: function (response) {
                // Handle the successful response
                console.log(response);

                // Remove the deleted row from the table
                //$('#' + taskId).remove();
            },
            error: function (xhr, status, error) {
                // Handle the error
                console.error(error);
            }
        });
    });

    // Add event listener to editable fields
    $('table').on('keyup focusout', '.editable', function (event) {
        var taskId = $(this).closest('tr').attr('id');
        var field = $(this).data('field');
        var currentValue = $(this).text();

        if (event.type === 'keyup' && event.keyCode !== 13) {
            return; // Skip if it's not the Enter key
        }

        var newValue = $(this).val();

        // Send an AJAX request to update the task
        $.ajax({
            url: '/tasks/update',
            type: 'POST',
            data: {
                taskId: taskId,
                field: field,
                value: newValue
            },
            success: function (response) {
                // Handle the successful response
                console.log(response);

                // Update the field with the new value
                $(this).text(newValue);
            }.bind(this),
            error: function (xhr, status, error) {
                // Handle the error
                console.error(error);

                // Revert the field back to the original value
                $(this).text(currentValue);
            }.bind(this)
        });
    });
});
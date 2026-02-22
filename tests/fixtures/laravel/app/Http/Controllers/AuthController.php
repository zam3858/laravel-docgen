<?php

namespace App\Http\Controllers;

class AuthController
{
    public function login()
    {
        return response()->json(['ok' => true]);
    }
}

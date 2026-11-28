<?php

namespace Tests\Feature;

use App\Models\ProgramStudi;
use App\Models\Mahasiswa;
use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class MahasiswaApiTest extends TestCase
{
    use RefreshDatabase;

    public function test_mahasiswa_requires_authentication(): void
    {
        $this->getJson('/api/mahasiswa')->assertUnauthorized();
    }

    public function test_authenticated_user_can_list_mahasiswa(): void
    {
        $this->withToken($this->tokenWithAbilities('mahasiswa:baca'))
            ->getJson('/api/mahasiswa')
            ->assertOk()
            ->assertJsonPath('meta.total', 0);
    }

    public function test_authorized_user_can_create_mahasiswa(): void
    {
        $programStudi = ProgramStudi::create([
            'kode' => 'TI',
            'nama' => 'Teknik Informatika',
            'jenjang' => 'S1',
        ]);

        $this->withToken($this->tokenWithAbilities('mahasiswa:baca', 'mahasiswa:tulis'))
            ->postJson('/api/mahasiswa', [
                'program_studi_id' => $programStudi->id,
                'nim' => 'H12345',
                'nama' => 'Mahasiswa Uji',
                'email' => 'mahasiswa@example.test',
                'angkatan' => 2024,
                'ipk' => 3.5,
                'aktif' => true,
            ])
            ->assertCreated()
            ->assertJsonPath('data.nim', 'H12345');

        $this->assertDatabaseHas('mahasiswas', [
            'nim' => 'H12345',
            'email' => 'mahasiswa@example.test',
        ]);
    }

    public function test_invalid_mahasiswa_input_returns_validation_errors(): void
    {
        $this->withToken($this->tokenWithAbilities('mahasiswa:baca', 'mahasiswa:tulis'))
            ->postJson('/api/mahasiswa', [])
            ->assertUnprocessable()
            ->assertJsonValidationErrors([
                'program_studi_id',
                'nim',
                'nama',
                'email',
                'angkatan',
            ]);
    }

    public function test_authorized_user_can_update_mahasiswa(): void
    {
        $programStudi = ProgramStudi::create([
            'kode' => 'TI',
            'nama' => 'Teknik Informatika',
            'jenjang' => 'S1',
        ]);

        $mahasiswa = Mahasiswa::create([
            'program_studi_id' => $programStudi->id,
            'nim' => 'H12345',
            'nama' => 'Nama Lama',
            'email' => 'lama@example.test',
            'angkatan' => 2024,
            'ipk' => 3.25,
            'aktif' => true,
        ]);

        $this->withToken($this->tokenWithAbilities('mahasiswa:baca', 'mahasiswa:tulis'))
            ->putJson('/api/mahasiswa/' . $mahasiswa->id, [
                'program_studi_id' => $programStudi->id,
                'nim' => 'H12345',
                'nama' => 'Nama Baru',
                'email' => 'baru@example.test',
                'angkatan' => 2026,
                'ipk' => 3.75,
            ])
            ->assertOk()
            ->assertJsonPath('data.nama', 'Nama Baru')
            ->assertJsonPath('data.email', 'baru@example.test')
            ->assertJsonPath('data.program_studi_id', $programStudi->id);

        $this->assertDatabaseHas('mahasiswas', [
            'id' => $mahasiswa->id,
            'nama' => 'Nama Baru',
            'email' => 'baru@example.test',
            'angkatan' => 2026,
            'ipk' => 3.75,
        ]);
    }

    public function test_read_only_token_cannot_create_matakuliah(): void
    {
        $this->withToken($this->tokenWithAbilities('mahasiswa:baca'))
            ->postJson('/api/matakuliah', [])
            ->assertForbidden();
    }

    public function test_login_is_rate_limited(): void
    {
        for ($attempt = 1; $attempt <= 6; $attempt++) {
            $response = $this->postJson('/api/auth/login', [
                'email' => 'unknown@example.test',
                'password' => 'incorrect-password',
            ]);

            if ($attempt < 6) {
                $response->assertUnprocessable();
            } else {
                $response->assertTooManyRequests();
            }
        }
    }

    private function tokenWithAbilities(string ...$abilities): string
    {
        $user = User::factory()->create(['peran' => 'admin']);

        return $user->createToken('test-token', $abilities)->plainTextToken;
    }
}

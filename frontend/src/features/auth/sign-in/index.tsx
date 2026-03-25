import { useSearch, Link } from '@tanstack/react-router'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { AuthLayout } from '../auth-layout'
import { UserAuthForm } from './components/user-auth-form'

export function SignIn() {
  const { redirect } = useSearch({ from: '/(auth)/sign-in' })

  return (
    <AuthLayout>
      <Card className='gap-4'>
        <CardHeader>
          <CardTitle className='text-lg tracking-tight'>Acesso ao Sistema</CardTitle>
          <CardDescription>
            Insira seu email e senha abaixo para acessar sua conta. <br />
            Ainda não tem cadastro?{' '}
            <Link
              to='/sign-up'
              className='hover:text-primary underline underline-offset-4'
            >
              Criar Conta
            </Link>
          </CardDescription>
        </CardHeader>
        <CardContent>
          <UserAuthForm redirectTo={redirect} />
        </CardContent>
        <CardFooter>
          <p className='text-muted-foreground px-8 text-center text-sm'>
            Ao entrar, você concorda com nossos{' '}
            <a
              href='/terms'
              className='hover:text-primary underline underline-offset-4'
            >
              Termos de Serviço
            </a>{' '}
            e{' '}
            <a
              href='/privacy'
              className='hover:text-primary underline underline-offset-4'
            >
              Política de Privacidade
            </a>
            .
          </p>
        </CardFooter>
      </Card>
    </AuthLayout>
  )
}
